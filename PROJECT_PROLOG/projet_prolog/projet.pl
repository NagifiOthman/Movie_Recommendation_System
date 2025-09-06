:- dynamic user/3, movie/2.
% K
min_liked(10).
% R
liked_th(3.5).
% N
number_of_rec(20).

read_users(Filename) :-
    csv_read_file(Filename, Data), assert_users(Data).
assert_users([]).
assert_users([row(U,_,_,_) | Rows]) :- \+number(U),!, assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U),\+user(U,_,_), liked_th(R), Rating>=R,!,assert(user(U,[M],[])), assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U),\+user(U,_,_), liked_th(R), Rating<R,!,assert(user(U,[],[M])), assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U), liked_th(R), Rating>=R, !, retract(user(U,Liked,NotLiked)), assert(user(U,[M|Liked],NotLiked)), assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U), liked_th(R), Rating<R, !, retract(user(U,Liked,NotLiked)), assert(user(U,Liked,[M|NotLiked])), assert_users(Rows).

read_movies(Filename) :-
    csv_read_file(Filename, Rows), assert_movies(Rows).

assert_movies([]).
assert_movies([row(M,_,_) | Rows]) :- \+number(M),!, assert_movies(Rows).
assert_movies([row(M,Title,_) | Rows]) :- number(M),!, assert(movie(M,Title)), assert_movies(Rows).

display_first_n(_, 0) :- !.
display_first_n([], _) :- !.
display_first_n([H|T], N) :-
    writeln(H), 
    N1 is N - 1,
    display_first_n(T, N1).

%! Calculate the similarity between two users
% User1: First user ID
% User2: Second user ID
% Sim: Similarity score
similarity(User1, User2, Sim) :-
    user(User1, Liked1, NotLiked1),
    user(User2, Liked2, NotLiked2),
    intersection(Liked1, Liked2, CommonLiked),
    length(CommonLiked, CommonLikedLen),
    intersection(NotLiked1, NotLiked2, CommonNotLiked),
    length(CommonNotLiked, CommonNotLikedLen),
    append(Liked1, NotLiked1, Movies1),
    append(Liked2, NotLiked2, Movies2),
    union(Movies1, Movies2, AllMovies),
    length(AllMovies, AllMoviesLen),
    Sim is (CommonLikedLen + CommonNotLikedLen) / AllMoviesLen.


%! Deteminese whether the user has seen and likes the movie
% Movie: Movie ID
% User: User ID
likes_movie(Movie, User) :-
    user(User, Liked, _),
    member(Movie, Liked).

%! Extracts the users in the list who liked the movie
% Movie: Movie ID
% Users: List of users
% UserWhoLiked: Subset of users that liked the movie
liked(Movie, Users, UserWhoLiked) :-
    include(likes_movie(Movie), Users, UserWhoLiked).

%! Calculates the probability that the user will like the movie
% User: User ID
% Movie: Movie ID
% Prob: Probability score
prob(User, Movie, Prob) :-
    findall(U, (user(U, _, _), U \= User), OtherUsers),
    liked(Movie, OtherUsers, UserWhoLiked),
    min_liked(K),
    length(UserWhoLiked, NUsers),
    (   
        NUsers < K
    ->  
        Prob is 0
    ;   
        maplist(similarity(User), UserWhoLiked, SimList),
        sum_list(SimList, SimSum),
        Prob is (SimSum / NUsers)
    ).


%! Deteminese whether the user has seen the movie
% Movie: Movie ID
% User: User ID
seen(User, Movie) :-
    user(User, Liked, NotLiked),
    (member(Movie, Liked); member(Movie, NotLiked)).

% Get the (Title, Prob) pair for this user and movie.
% User: User ID
% Movie: Movie ID
% Title: Movie title
% Prob: Recommendation probability
get_rec(User, Movie, (Title, Prob)) :-
    \+ seen(User, Movie),
    movie(Movie, Title),
    prob(User, Movie, Prob).

%! Generate the (Title, Prob) pairs for the movies and the user
% User: User ID
% Movies: List of movie IDs
% Recs: List of (Title, Prob) pairs representing recommendations
prob_movies(User, Movies, Recs) :-
    findall(Rec, (member(M, Movies), get_rec(User, M, Rec)), Recs).


recommendations(User) :-
    setof(M,L^movie(M,L),Ms),   % generate list of all movie 
	prob_movies(User,Ms,Rec),   % compute probabilities for all movies 
	sort(2,@>=,Rec,Rec_Sorted), % sort by descending probabilities
	number_of_rec(N),
    display_first_n(Rec_Sorted,N). % display the result

init :- read_users('ml-latest-small/ratings.csv'), read_movies('ml-latest-small/movies.csv').
test(1):- similarity(33,88,S1), 291 is truncate(S1 * 10000),similarity(44,55,S2), 138 is truncate(S2 * 10000).
test(2):- prob(44,1080,P1), 122 is truncate(P1 * 10000), prob(44,1050,P2), 0 is truncate(P2).
test(3):- liked(1080, [28, 30, 32, 40, 45, 48, 49, 50], [28, 45, 50]).
test(4):- seen(32, 1080), \+seen(44, 1080).
test(5):- prob_movies(44,[1010, 1050, 1080, 2000],Rs), length(Rs,4), display(Rs).


                                   