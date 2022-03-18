sudo yum install -y https://yum.postgresql.org/14/redhat/rhel-7-x86_64/postgresql14-libs-14.2-1PGDG.rhel7.x86_64.rpm
sudo yum install -y https://yum.postgresql.org/14/redhat/rhel-7-x86_64/postgresql14-14.2-1PGDG.rhel7.x86_64.rpm
sudo yum install -y https://yum.postgresql.org/14/redhat/rhel-7-x86_64/postgresql14-server-14.2-1PGDG.rhel7.x86_64.rpm
sudo /usr/pgsql-14/bin/postgresql-14-setup initdb
sudo systemctl enable postgresql-14
sudo systemctl start postgresql-14

sudo sed -i -e "/listen_addresses/s/localhost/*/" /var/lib/pgsql/14/data/postgresql.conf
sudo sed -i -e "/#listen_addresses/s/#listen_addresses/listen_addresses/" /var/lib/pgsql/14/data/postgresql.conf
echo 'host all all 0.0.0.0/0 md5' | sudo tee -a /var/lib/pgsql/14/data/pg_hba.conf

echo "CREATE DATABASE mydatabase WITH ENCODING 'UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8';
GRANT TEMP ON DATABASE mydatabase TO postgres;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO postgres;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO postgres;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO postgres;
create user admin with superuser encrypted password 'admin';
grant all privileges on database mydatabase to admin;" | sudo -u postgres psql

echo "\c mydatabase
CREATE TABLE messages ( uuid CHAR (36) PRIMARY KEY,
timestamp TIMESTAMP WITH TIME ZONE NULL,
author VARCHAR (64) NOT NULL,
message VARCHAR (1024) NOT NULL,
likes INT NOT NULL);" | sudo -u postgres psql

echo "\c mydatabase
CREATE TABLE updated_messages ( uuid CHAR (36) PRIMARY KEY,
timestamp TIMESTAMP WITH TIME ZONE NULL,
author VARCHAR (64) NOT NULL,
message VARCHAR (1024) NOT NULL,
likes INT NOT NULL,
is_deleted BOOLEAN NOT NULL );" | sudo -u postgres psql

sudo reboot
