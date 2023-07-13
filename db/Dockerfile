# Use the MongoDB Docker image
FROM mongo

# Copy the JSON files into the Docker image.
COPY data/users.json /users.json
COPY data/interviews.json /interviews.json
COPY data/comments.json /comments.json

# Import data and run MongoDB
CMD mongod --fork --logpath /var/log/mongod.log --bind_ip_all; \
    mongoimport --host localhost --db robinhood --collection users --type json --file /users.json --jsonArray; \
    mongoimport --host localhost --db robinhood --collection interviews --type json --file /interviews.json --jsonArray; \
    mongoimport --host localhost --db robinhood --collection comments --type json --file /comments.json --jsonArray; \
    tail -f /var/log/mongod.log
