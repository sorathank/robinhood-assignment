# Use the MongoDB Docker image
FROM mongo:4.4.6

# Copy the JSON files into the Docker image.
COPY data.json /data.json
COPY interviews.json /interviews.json
COPY comments.json /comments.json

# Make sure the MongoDB server starts and then import the data.
CMD mongod --fork --logpath /var/log/mongod.log --bind_ip_all && \
    mongoimport --host localhost --db robinhood --collection users --type json --file /data/users.json --jsonArray && \
    mongoimport --host localhost --db robinhood --collection interviews --type json --file /data/interviews.json --jsonArray && \
    mongoimport --host localhost --db robinhood --collection comments --type json --file /data/comments.json --jsonArray
