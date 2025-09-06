import java.util.ArrayList;
import java.util.List;

public class User {
    private int userId;
    private List<Movie> likedMovies;
    private List<Movie> dislikedMovies;

    public User(int userId) {
        this.userId = userId;
        this.likedMovies = new ArrayList<>();
        this.dislikedMovies = new ArrayList<>();
    }

    public int getUserId() {
        return userId;
    }

    public List<Movie> getLikedMovies() {
        return likedMovies;
    }

    public List<Movie> getDislikedMovies() {
        return dislikedMovies;
    }

    public void likeMovie(Movie m) {
        likedMovies.add(m);
    }

    public void dislikeMovie(Movie m) {
        dislikedMovies.add(m);
    }

    @Override
    public String toString() {
        return "User{" +
                "userId=" + userId +
                ", likedMovies=" + likedMovies.size() +
                ", dislikedMovies=" + dislikedMovies.size() +
                '}';
    }
}
