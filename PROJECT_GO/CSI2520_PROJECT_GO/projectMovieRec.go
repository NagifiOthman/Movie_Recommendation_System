// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// movies with rating greater or equal are considered 'liked'
const iLiked float64 = 3.5

// Define the Recommendation type
type Recommendation struct {
	userID     int     // recommendation for this user
	movieID    int     // recommended movie ID
	movieTitle string  // recommended movie title
	score      float32 // probability that the user will like this movie
	nUsers     int     // number of users who likes this movie
}

// get the probability that this user will like this movie
func (r Recommendation) getProbLike() float32 {
	return r.score / (float32)(r.nUsers)
}

// Define the User type
// and its list of liked items
type User struct {
	userID   int
	liked    []int // list of movies with ratings >= iLiked
	notLiked []int // list of movies with ratings < iLiked
}

func (u User) getUser() int {
	return u.userID
}

func (u *User) setUser(id int) {
	u.userID = id
}

func (u *User) addLiked(id int) {
	u.liked = append(u.liked, id)
}

func (u *User) addNotLiked(id int) {
	u.notLiked = append(u.notLiked, id)
}

// intersectionCount returns the number of common elements between two slices.
func intersectionCount(a, b []int) int {
	count := 0
	set := make(map[int]bool)
	for _, x := range a {
		set[x] = true
	}
	for _, x := range b {
		if set[x] {
			count++
		}
	}
	return count
}

// unionCountOfTwoUsers returns the size of the union of the movies (liked and notLiked)
// that appear in either user's history.
func unionCountOfTwoUsers(u1, u2 *User) int {
	set := make(map[int]bool)
	for _, x := range u1.liked {
		set[x] = true
	}
	for _, x := range u1.notLiked {
		set[x] = true
	}
	for _, x := range u2.liked {
		set[x] = true
	}
	for _, x := range u2.notLiked {
		set[x] = true
	}
	return len(set)
}

// computeSimilarity computes the Jaccard similarity between two users,
// based on their liked and notLiked movies.
func computeSimilarity(u1, u2 *User) float32 {
	interLikes := intersectionCount(u1.liked, u2.liked)
	interDislikes := intersectionCount(u1.notLiked, u2.notLiked)
	union := unionCountOfTwoUsers(u1, u2)
	if union == 0 {
		return 0
	}
	return float32(interLikes+interDislikes) / float32(union)
}

func computeMovieLikes(users map[int]*User) map[int][]int {
	movieLikes := make(map[int][]int)
	for uid, user := range users {
		for _, movieID := range user.liked {
			movieLikes[movieID] = append(movieLikes[movieID], uid)
		}
	}
	return movieLikes
}

// Function to read the ratings CSV file and process each row.
// The output is a map in which user ID is used as key
func readRatingsCSV(fileName string) (map[int]*User, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader.
	reader := csv.NewReader(file)

	// Read first line and skip
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// creates the map
	users := make(map[int]*User, 1000)

	// Read all records from the CSV.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate over each record and convert the strings into integers or float.
	for _, record := range records {
		if len(record) != 4 {
			return nil, fmt.Errorf("each line must contain exactly 4 integers, but found %d", len(record))
		}

		// Parse user ID integer
		uID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to userID integer: %v", record[0], err)
		}

		// Parse movie ID integer
		mID, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to movieID integer: %v", record[1], err)
		}

		// Parse rating float
		r, err := strconv.ParseFloat(record[2], 64)
		if err != nil {

			return nil, fmt.Errorf("error converting '%s' to rating: %v", record[2], err)
		}

		// checks if it is a new user
		u, ok := users[uID]
		if !ok {

			u = &User{uID, nil, nil}
			users[uID] = u
		}

		// ad movie in user list
		if r >= iLiked {

			u.addLiked(mID)

		} else {

			u.addNotLiked(mID)
		}
	}

	return users, nil
}

// Function to read the movies CSV file and process each row.
// The output is a map in which user ID is used as key
func readMoviesCSV(fileName string) (map[int]string, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader.
	reader := csv.NewReader(file)

	// Read first line and skip
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// creates the map
	movies := make(map[int]string, 1000)

	// Read all records from the CSV.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate over each record and convert the strings into integers or float.
	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("each line must contain exactly 3 entries, but found %d", len(record))
		}

		// Parse movie ID integer
		mID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to movieID integer: %v", record[1], err)
		}

		// record 1 is the title
		movies[mID] = record[1]
	}

	return movies, nil
}

// checks if value is in the set
func member(value int, set []int) bool {

	for _, v := range set {
		if value == v {

			return true
		}
	}

	return false
}

// generator producing Recommendation instances from movie list
func generateMovieRec(wg *sync.WaitGroup, stop <-chan bool, userID int, titles map[int]string) <-chan Recommendation {

	outputStream := make(chan Recommendation)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer close(outputStream)
		// defer fmt.Println("\nFin de generateMovieRec...")
		for k, v := range titles {
			select {
			case <-stop:
				return
			case outputStream <- Recommendation{userID, k, v, 0.0, 0}:
			}
		}
	}()

	return outputStream
}

func computeLikeCount(users map[int]*User) map[int]int {
	likeCount := make(map[int]int)

	for _, u := range users {
		for _, movieID := range u.liked {
			likeCount[movieID]++
		}
	}
	return likeCount
}

func filterSeenMovies(
	wg *sync.WaitGroup,
	stop <-chan bool,
	in <-chan Recommendation,
	user *User,
) <-chan Recommendation {
	out := make(chan Recommendation)
	// Do NOT call wg.Add(1) here; main already did it.
	go func() {
		defer wg.Done()
		defer close(out)
		for rec := range in {
			select {
			case <-stop:
				return
			default:
				// Pass this recommendation along ONLY IF the user hasn’t seen the movie.
				if !member(rec.movieID, user.liked) && !member(rec.movieID, user.notLiked) {
					out <- rec
				}
			}
		}
	}()
	return out
}

func filterMinLikes(
	wg *sync.WaitGroup,
	stop <-chan bool,
	in <-chan Recommendation,
	likeCount map[int]int,
	K int,
) <-chan Recommendation {
	out := make(chan Recommendation)
	// Do NOT call wg.Add(1) here; main already did it.
	go func() {
		defer wg.Done()
		defer close(out)
		for rec := range in {
			select {
			case <-stop:
				return
			default:
				if likeCount[rec.movieID] >= K {
					out <- rec
				}
			}
		}
	}()
	return out
}

func computeScore(rec Recommendation, currentUser *User, users map[int]*User, movieLikes map[int][]int) Recommendation {
	likedUsers, ok := movieLikes[rec.movieID]
	if !ok || len(likedUsers) == 0 {
		// If no one liked the movie, set score to 0.
		rec.score = 0
		rec.nUsers = 0
		return rec
	}

	var totalSim float32 = 0.0
	for _, vID := range likedUsers {
		// Skip if the user is the current user.
		if vID == currentUser.userID {
			continue
		}
		vUser, exists := users[vID]
		if !exists {
			continue
		}
		sim := computeSimilarity(currentUser, vUser)
		totalSim += sim
	}
	rec.nUsers = len(likedUsers)
	// Average similarity is the score.
	rec.score = totalSim / float32(len(likedUsers))
	return rec
}

func scoreStage(
	wg *sync.WaitGroup,
	stop <-chan bool,
	in <-chan Recommendation,
	out chan<- Recommendation,
	currentUser *User,
	users map[int]*User,
	movieLikes map[int][]int,
) {
	defer wg.Done()
	for rec := range in {
		select {
		case <-stop:
			return
		default:
			scored := computeScore(rec, currentUser, users, movieLikes)
			out <- scored
		}
	}
}

func runParallelScoring(
	stop <-chan bool,
	in <-chan Recommendation,
	currentUser *User,
	users map[int]*User,
	movieLikes map[int][]int,
) <-chan Recommendation {
	out := make(chan Recommendation)
	var wgLocal sync.WaitGroup
	wgLocal.Add(2) // Two parallel scoring goroutines

	go scoreStage(&wgLocal, stop, in, out, currentUser, users, movieLikes)
	go scoreStage(&wgLocal, stop, in, out, currentUser, users, movieLikes)

	go func() {
		wgLocal.Wait()
		close(out)
	}()
	return out
}

func collectAndShowResults(in <-chan Recommendation, N int) {
	var results []Recommendation
	for r := range in {
		results = append(results, r)
	}
	// sort results by score in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	// Print the top N recommendations in the desired format
	for i := 0; i < N && i < len(results); i++ {
		rec := results[i]
		fmt.Printf("%s at %.4f [ %d]\n", rec.movieTitle, rec.score, rec.nUsers)
	}
}

func main() {
	// Which user do we recommend for?
	var currentUserID int
	fmt.Println("Enter user ID :")
	fmt.Scanf("%d", &currentUserID)

	// Read movies and ratings CSVs
	titles, err := readMoviesCSV("movies.csv")
	if err != nil {
		log.Fatal(err)
	}
	ratings, err := readRatingsCSV("ratings.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Precompute how many users like each movie (if needed) and compute movieLikes.
	likeCountMap := computeLikeCount(ratings) // you may keep this if used elsewhere
	movieLikes := computeMovieLikes(ratings)

	// Config values for filtering and output
	K := 10 // filter out movies liked by fewer than 10 users
	N := 20 // display top 21 recommendations

	// Confirm this user exists in the ratings map
	currentUser, ok := ratings[currentUserID]
	if !ok {
		fmt.Printf("User %d not found in ratings data.\n", currentUserID)
		return
	}

	// Measure execution time
	start := time.Now()

	// Channel and goroutine management for stages 1-3 using a global WaitGroup.
	stop := make(chan bool)
	var wg sync.WaitGroup

	// Stage 1: generate a Recommendation for every movie.
	wg.Add(1)
	genChan := generateMovieRec(&wg, stop, currentUserID, titles)

	// Stage 2: remove movies the user has already seen.
	wg.Add(1)
	unseenChan := filterSeenMovies(&wg, stop, genChan, currentUser)

	// Stage 3: only keep movies liked by >= K users.
	wg.Add(1)
	filteredChan := filterMinLikes(&wg, stop, unseenChan, likeCountMap, K)

	// Stage 4: run two parallel scoring goroutines.
	scoredChan := runParallelScoring(stop, filteredChan, currentUser, ratings, movieLikes)

	// Final: collect, sort, and print top N recommendations.
	fmt.Printf("Recommendations for user # %d:\n", currentUserID)
	collectAndShowResults(scoredChan, N)

	// Wait for the pipeline stages (generator and filters) to finish.
	wg.Wait()

	// Print total execution time.
	end := time.Now()
	fmt.Printf("\nExecution time: %v\n", end.Sub(start))
}
