public class Recommendation {
    private User user;
    private Movie movie;
    private double score;    // probability that 'user' will like 'movie'
    private int nUsers;      // number of users who liked this movie

    public Recommendation(User user, Movie movie, double score, int nUsers) {
        this.user = user;
        this.movie = movie;
        this.score = score;
        this.nUsers = nUsers;
    }

    public User getUser() {
        return user;
    }

    public Movie getMovie() {
        return movie;
    }

    public double getScore() {
        return score;
    }

    public int getNUsers() {
        return nUsers;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public void setMovie(Movie movie) {
        this.movie = movie;
    }

    public void setScore(double score) {
        this.score = score;
    }

    public void setNUsers(int nUsers) {
        this.nUsers = nUsers;
    }



    @Override
    public String toString() {
        return "Recommendation{" +
                "userId=" + user.getUserId() +
                ", movieId=" + movie.getMovieId() +
                ", probability=" + score +
                ", nUsersWhoLiked=" + nUsers +
                '}';
    }
}
