var carousel = $('.carousel-inner');
$.ajax({
    url: window.origin + "/api/category/all",
    method: "POST",
    success: function(data){
        var movie = JSON.parse(data);
        var movieCards = [];
        var row = $('.movie-row');
     
        for(var i = 0; i < movie.length; i++) {
            movieCards[i] = ` <div class="col-md-2">
            <div class="card m-0 p-0">
                <img class="card-img" src="${movie[i].MovieInfo.Poster}" alt="Movie Poster">

                <div class="card-img-overlay">
                    <div class="d-flex flex-column align-items-center">
                        <div style="margin-top: 50%;"></div>
                        <i class="fab fa-imdb imdb-logo"></i>
                        <span class="overlay-text">${movie[i].MovieInfo.Ratings[0].Value}</span>
                    </div>
                </div>
            </div>
        </div>`
        }
        addMovies(movieCards)
    },
    error: function(x, y, z) {
        console.log(x, y ,z)
    }

});

// add every 6 movies to a carousel item row
function addMovies(movies){
    var maxCards = 6; //max card per row
    var rowsQnt = Math.ceil(movies.length);
    var rows = [];

    if(rowsQnt > 1){

    }

    for(var i = 0; i < rowsQnt; i++){
        rows[i] = `<div class="carousel-item active">
        <div class="row movie-row">`;

        for(var movieIndex = 0; movieIndex < maxCards; movieIndex++){
            if(movies[movieIndex] != null)
                rows[i] += movies[movieIndex];

        }
        rows[i] += `</div></div>`;

    }
    carousel.append(rows);
}

