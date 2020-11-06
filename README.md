### Mini_URL

#### Running the Application

To run the application, run the command `docker-compose up -d`. The application takes around 10 seconds waiting for the database to be ready, please check its container logs to make sure. To see what API endpoints are offered by the application, visit `http://localhost:4000/swagger/index.html`

#### Shortening a URL

To shorten a URL, send a `POST http://localhost:4000/`, the URL you get back is the shortened URL

#### Using a Shortened URL

To use a shortened URL, simply type `http://localhost:4000/{shortURL}`