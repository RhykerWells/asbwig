<p align="center">
  <a href="https://github.com/RhykerWells/asbwig">
    <picture>
      <img src="./frontend/static/img/avatar.png" height="128">
    </picture>
    <h1 align="center">ASBWIG</h1>
    <h1 align="center">
      <span style="font-weight: bold;">A</span>nother <span style="font-weight: bold;">S</span>hitty</span> <span style="font-weight: bold;">B</span>ot</span> <span style="font-weight: bold;">W</span>ritten</span> <span style="font-weight: bold;">I</span>n</span> <span style="font-weight: bold;">G</span>o</span>
    </h1>
  </a>
</p>
ASBWIG is a bot I've decided to write to give myself a reason to learn Go and how the language works.
It by all means will not be perfect as I am teaching myself this language as I progress in this.
A lot of the inspiration for command structure and the like comes from <a href="https://github.com/botlabs-gg/yagpdb">YAGPDB</a>

# Selfhosting
## Standalone
[Install Golang](https://go.dev/doc/install)

Install Postgres</br>
```
sudo apt update
sudo apt install postgresql
```

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
git clone https://github.com/RhykerWells/asbwig
cd asbwig/cmd/asbwig
go build
```

Once it has finished compiling. Run the binary with:</br>
`./asbwig`

## Docker

Update your system and install git
```
sudo apt update
sudo apt install git
```
Clone the repository
```
git clone https://github.com/RhykerWells/asbwig
cd asbwig/docker
```
Copy the environment variable files and edit where applicable
```
cp app.example, app.env
cp db.example, db.env
```
Add your docker image to the compose file
`docker compose up -d`