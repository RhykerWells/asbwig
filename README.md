<p align="center">
  <a href="https://github.com/ranger-4297/asbwig">
    <picture>
      <img src="./img/avatar.png" height="128">
    </picture>
    <div style="text-align:center">
      <svg xmlns="http://www.w3.org/2000/svg" height="0">
      <style>
        .first-letter {
          font-weight: bold;
        }
      </style>
      <h1>ASBWIG</h1>
      <h1 class="name"><span class="first-letter">A</span>nother <span class="first-letter">S</span>hitty</span> <span class="first-letter">B</span>ot</span> <span class="first-letter">W</span>ritten</span> <span class="first-letter">I</span>n</span> <span class="first-letter">G</span>o</span></h1>
      </svg>
    </div>
  </a>
</p>
</style>
ASBWIG is a bot I've decided to write to give myself a reason and an idea on how the Go language works.
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
git clone https://github.com/Ranger-4297/asbwig
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
git clone https://github.com/Ranger-4297/asbwig
cd asbwig/docker
```
Copy the environment variable files and edit where applicable
```
cp app.example, app.env
cp db.example, db.env
```
Add your docker image to the compose file
`docker compose up -d`