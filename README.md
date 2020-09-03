# SpotifyParty Bot

SpotifyParty Bot lets you listen to music with your friends in real time through discord and spotify.


[Invite to your server](https://discord.com/api/oauth2/authorize?client_id=747675945644851221&permissions=19520&scope=bot)

## How it works

First, you need to create a **Party** on your server and ask your friends to join. Each user will need to connect with Spotify in a link sent to the DM by the **SpotifyParty Bot**. After this, all party participants can add music to play in the **Party** for everyone.

When a song plays, it will play synchronously for everyone just like in a party :partying_face:.

> Obs: The bot will respect the sequence of one user song at a time, so everyone can enjoy it equally.

## List of Commands

Command | Description
:---: | ---
**sp.join** | Join a party. If no party exists in the server it will also create.
**sp.add** | Add a new music to the users queue. You can check the musics added on Spotify.
**sp.list** | List all users participating in the party.
**sp.sync** | Sync music and time with the party on Spotify.
**sp.leave** | Leave the party.
**sp.kill** | Kill the party.
**sp.help** | List of Commands.


## How to Contribute

Pull requests and issues are welcome! Feel free to make suggestions, report bugs and help with new features :heart:.

### Development 

First of all, create a new file `.env` with the same fields as the `.env.example` file. You will use this file to fill with Spotify and Discord data.

#### Spotify

Access the [Spotify Dashboard](https://developer.spotify.com/dashboard/) and create a new Application. Follow these steps:

1. Put the *Client ID* and *Client Secret* ID in the `.env`
2. Go to **Edit Settings** on App's page
3. Add `http://localhost:8080/auth` to **Redirect URLs**

#### Discord

Access [Discord Developer Portal](https://discord.com/developers/applications) and create a new Application, also create a bot inside of this application.

1. Go in **Bot**
2. Copy the *Token* and add it to the `.env`
3. Go to **OAuth**, Check *bot* as scope and add the follow permissions: See Channels., Send Messages, Manage Messages, Embed Links, Read Message History, Add reactions
4. Copy the invitation link and add to your server

#### Run the bot

```sh
go run *.go
```

## License

Licensed under the [MIT](./LICENSE) license.