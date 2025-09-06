// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca

public class Movie {
	
	private int movieID;
	private String title;
	
	// constructs a movie
    public Movie(int id, String title) {
	
		movieID= id;
		this.title= title;		
	}
	
	// gets the ID
	public int getMovieId() {
	
		return movieID;
	}

	// get the movie title
	public String getTitle() {
	
		return title;
	}
	
	// string representation
	public String toString() {
      
       return "["+movieID+"]:"+title;	  
	}
}