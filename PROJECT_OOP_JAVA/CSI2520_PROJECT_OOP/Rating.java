public class Rating {
    private int userID;
    private int movieID;
    private float rating;

    private float timestamp;

    public Rating(int userID, int movieID, float rating, float timestamp) {
        this.userID = userID;
        this.movieID = movieID;
        this.rating = rating;
        this.timestamp = timestamp;
    }

    public float getTimestamp() {
        return timestamp;
    }


    public int getUserID() {
        return userID;
    }

    public int getMovieID() {
        return movieID;
    }

    public float getRating() {
        return rating;
    }



    @Override
    public String toString() {
        return "Rating[userID=" + userID + ", movieID=" + movieID
                + ", rating=" + rating + "]";
    }
}
