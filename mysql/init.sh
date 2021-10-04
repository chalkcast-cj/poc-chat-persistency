mysql -u root -ppassword -e "create database if not exists chat;"

mysql -u root -ppassword -e "CREATE TABLE chat.messages (
   message_id VARCHAR(30) NOT NULL,
   user_id  INT NOT NULL,
   room_id INT NOT NULL,
   message TEXT NOT NULL,
   PRIMARY KEY ( message_id )
);"
