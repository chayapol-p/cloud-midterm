sudo yum install -y https://yum.postgresql.org/14/redhat/rhel-7-x86_64/postgresql14-libs-14.2-1PGDG.rhel7.x86_64.rpm
sudo yum install -y https://yum.postgresql.org/14/redhat/rhel-7-x86_64/postgresql14-14.2-1PGDG.rhel7.x86_64.rpm
sudo yum install -y https://yum.postgresql.org/14/redhat/rhel-7-x86_64/postgresql14-server-14.2-1PGDG.rhel7.x86_64.rpm
sudo /usr/pgsql-14/bin/postgresql-14-setup initdb
sudo systemctl enable postgresql-14
sudo systemctl start postgresql-14

sudo -u postgres createuser --superuser ec2-user
CREATE DATABASE `myDatabase` WITH ENCODING 'UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8';
GRANT TEMP ON DATABASE `myDatabase` TO `postgres`;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO `postgres`;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO `postgres`;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO `postgres`;