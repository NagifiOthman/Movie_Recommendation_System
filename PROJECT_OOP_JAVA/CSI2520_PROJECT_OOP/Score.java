public class Score {
    Movie movie;
    double score;
    int userCount;  // New attribute to store the number of users who liked this movie

    Score(Movie movie, double score, int userCount) {
        this.movie = movie;
        this.score = score;
        this.userCount = userCount;
    }

    public Movie getMovie() {
        return movie;
    }

    public double getScore() {
        return score;
    }

    public int getUserCount() {
        return userCount;
    }
}
