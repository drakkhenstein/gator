This Aggrogator Program requires the instellation of Postgres and Go. 
Go is the Programming Language, its excutable only after compiling using: go build or go install. To test out the Program before compiling use: go run . at the root of the repo directory.  If you want to play around with modding the Program use: go run . agg to get it running in the Terminal, you need to open a second Terminal to type commands to use for it. 

Postgres is the SQL database package for this Program, it's used to generate code migrations and create/modify Tables.  It also runs queries in SQL to get specific data from the database tables.  The queries really make the Program commands simple. There are several commands you can use while running the Program, they are: register login users reset agg addfeed feeds follow following unfollow browse  each of these is using the queries code from the database to excute.

To install the gator program simply use: go install from the Command Line Interface, CLI, from the root of the repo directory.  From then on just use: gator (and a command) to run the Program.  

To setup the config file which sets up the .json code for the database's user/password/location permissions.  Simply use: touch .gatorconfig.json to create the gator config file in the root directory then add in the file: {"db_url":"postgres://password:username@localhostname:port/gator?sslmode=disable","current_user_name":username} with the db_url filled in with what your database's hostname:port:username:password when you run the Program the database logs you in with this config file.

The commands of the Program are used like this:

register - creates a new user for the database

login - switches to a new user

users - sets the user

reset - resets the users in the database

agg - collects the feeds/posts of users

addfeed - adds a new feed to the user

feeds - returns all the feeds of the user

follow - follows the feed of the user

following - lists the feeds followed by a user

unfollow - unfollows a feed of a user

browse - lists the posts of followed feeds