# Selfhosting
## Standalone
[Install Golang](https://go.dev/doc/install)

Install Postgres</br>
`sudo apt update`</br>
`sudo apt install postgresql`

Configure Postgres</br>
`sudo -u postgres psql`
```
CREATE DATABASE asbwig;
create user asbwig with encrypted password 'password';
grant all privileges on database asbwig to asbwig;
\c asbwig
grant usage, create on schema public to asbwig;
\q
```


Add your environment variables to your `~/.profile`</br>
`ASBWIG_TOKEN` - Your bot token. NOT prefixed with "Bot"</br>
`ASBWIG_PGUSERNAME` - The user in postgres you created</br>
`ASBWIG_PASSWORD` - The password you set in postgres

Prefix each variable with `export`:</br>
`export ASBWIG_TOKEN="tokenxxxx"`

Downloading and installing
```
sudo apt update
sudo apt install git
git clone https://github.com/Ranger-4297/ASBWIG
cd ASBWIG/cmd/ASBWIG
go build
```

Once it has finished compiling. Run the binary with:</br>
`./ASBWIG`

## Docker
