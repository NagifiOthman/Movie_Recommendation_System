// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca
import java.io.*;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.Set;

public class RecommendationEngine {

    public ArrayList<Movie> movies;
    private ArrayList<Rating> ratings = new ArrayList<>();
    private final double R = 3.5;
    private final int K = 10;
    private final int N = 20;

    // Constructs a recommendation engine from files
    public RecommendationEngine(String movieFile, String ratingsFile) throws IOException {
        movies = new ArrayList<>();
        readMovies(movieFile);
    }

    // Reads the Movie CSV file and populates the list of Movies
    public void readMovies(String csvFile) throws IOException {
        String line;
        String delimiter = ",";

        BufferedReader br = new BufferedReader(new FileReader(csvFile));
        br.readLine(); // Skip header

        while ((line = br.readLine()) != null && line.length() > 0) {
            String[] parts = line.split(delimiter);
            if (parts.length < 2) continue;

            int movieID = Integer.parseInt(parts[0]);
            String title = parts[1];

            movies.add(new Movie(movieID, title));
        }
        br.close();
    }

    // Reads the Ratings CSV file and returns a list of Ratings
    public ArrayList<Rating> readRatings(String csvFile) throws IOException {
        ArrayList<Rating> ratingList = new ArrayList<>();
        String line;

        BufferedReader br = new BufferedReader(new FileReader(csvFile));
        br.readLine(); // Skip header

        while ((line = br.readLine()) != null) {
            if (line.trim().isEmpty()) continue;

            String[] parts = line.split(",");
            if (parts.length < 4) continue;

            int userID = Integer.parseInt(parts[0]);
            int movieID = Integer.parseInt(parts[1]);
            float ratingValue = Float.parseFloat(parts[2]);
            float timeStamp = Float.parseFloat(parts[3]);

            ratingList.add(new Rating(userID, movieID, ratingValue, timeStamp));
        }
        br.close();
        return ratingList;
    }

    // Returns a list of movies watched by the user
    public ArrayList<Movie> moviesWatchedByUser(int userId, String csvFile) throws IOException {
        ArrayList<Movie> watched = new ArrayList<>();
        ArrayList<Rating> tempRatings = readRatings(csvFile);
        ArrayList<Integer> watchedMovieIDs = new ArrayList<>();

        for (Rating r : tempRatings) {
            if (r.getUserID() == userId) {
                watchedMovieIDs.add(r.getMovieID());
            }
        }

        for (Integer movieID : watchedMovieIDs) {
            Movie m = findMovieByID(movieID);
            if (m != null) {
                watched.add(m);
            }
        }
        return watched;
    }

    // Finds a movie by its ID
    private Movie findMovieByID(int movieID) {
        for (Movie m : movies) {
            if (m.getMovieId() == movieID) {
                return m;
            }
        }
        return null;
    }

    // Checks if a movie is liked by at least K users
    public boolean movieLiked(Movie movie, int k) {
        int count = 0;
        for (Rating r : ratings) {
            if (r.getMovieID() == movie.getMovieId() && r.getRating() >= R) {
                count++;
            }
            if (count >= k) {
                return true;
            }
        }
        return false;
    }

    // Generates recommendations for a specific user
    public ArrayList<Score> generateRecommendation(int userId, String csvfile) throws IOException {
        ArrayList<Movie> usersMovies = moviesWatchedByUser(userId, csvfile);
        double[] scoreUM = new double[movies.size()];
        int[] countLM = new int[movies.size()];  // Track how many users liked each movie
    
        for (int i = 0; i < movies.size(); i++) {
            scoreUM[i] = 0.0;
            countLM[i] = 0;
        }
    
        for (int i = 0; i < movies.size(); i++) {
            Movie m = movies.get(i);
    
            if (usersMovies.contains(m)) continue;
    
            if (movieLiked(m, K)) {
                for (Rating r : ratings) {
                    int v = r.getUserID();
                    if (v == userId) continue;
    
                    if (r.getMovieID() == m.getMovieId() && r.getRating() >= R) {
                        double sUV = computeSimilarity(userId, v);
                        scoreUM[i] += sUV;
                        countLM[i]++;  // Increment user count for this movie
                    }
                }
            }
        }
    
        ArrayList<Score> results = new ArrayList<>();
        for (int i = 0; i < movies.size(); i++) {
            if (countLM[i] > 0) {
                double pUM = scoreUM[i] / countLM[i];
                results.add(new Score(movies.get(i), pUM, countLM[i]));  // Include user count
            }
        }
    
        results.sort((a, b) -> Double.compare(b.getScore(), a.getScore()));
    
        return new ArrayList<>(results.subList(0, Math.min(N, results.size())));
    }
    

    // Computes similarity between two users using Jaccard Index
    private double computeSimilarity(int userId, int userVId) {
        Set<Integer> LU = new HashSet<>();
        Set<Integer> DU = new HashSet<>();
        for (Rating r : ratings) {
            if (r.getUserID() == userId) {
                if (r.getRating() >= R) {
                    LU.add(r.getMovieID());
                } else {
                    DU.add(r.getMovieID());
                }
            }
        }

        Set<Integer> LV = new HashSet<>();
        Set<Integer> DV = new HashSet<>();
        for (Rating r : ratings) {
            if (r.getUserID() == userVId) {
                if (r.getRating() >= R) {
                    LV.add(r.getMovieID());
                } else {
                    DV.add(r.getMovieID());
                }
            }
        }

        Set<Integer> commonLiked = new HashSet<>(LU);
        commonLiked.retainAll(LV);
        int numCommonLiked = commonLiked.size();

        Set<Integer> commonDisliked = new HashSet<>(DU);
        commonDisliked.retainAll(DV);
        int numCommonDisliked = commonDisliked.size();

        Set<Integer> unionAll = new HashSet<>(LU);
        unionAll.addAll(LV);
        unionAll.addAll(DU);
        unionAll.addAll(DV);
        int numUnion = unionAll.size();

        if (numUnion == 0) {
            return 0.0;
        }
        return (numCommonLiked + numCommonDisliked) / (double) numUnion;
    }

    // MAIN method with dynamic user ID support
    public static void main(String[] args) {
        if (args.length < 3) {
            System.err.println("Usage: java RecommendationEngine <userId> <movies.csv> <ratings.csv>");
            
            return;
        }
    
        try {
            int userId = Integer.parseInt(args[0]);
            String moviesFile = args[1];
            String ratingsFile = args[2];
            
            RecommendationEngine rec = new RecommendationEngine(moviesFile, ratingsFile);
            rec.ratings = rec.readRatings(ratingsFile);
    
            ArrayList<Score> recommended = rec.generateRecommendation(userId, ratingsFile);
    
            System.out.println( "Top "+rec.N+ " Recommendations for user # " + userId + ":");
            for (Score s : recommended) {
                System.out.println("movie: "+s.getMovie().getTitle()+"at: "+s.getScore()+" "+"["+ s.getUserCount()+"]");
            }
    
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
        }
    }
    
}
