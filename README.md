# dlfm
## This app will update your Discord status with information about song that you scrobble now

<details><summary><b>How to install</b></summary>

You can download a compiled version from [releases pages](https://github.com/dikey0ficial/dlfm/releases)

If you have go installed, you can only write this in console:
`go install github.com/dikey0ficia/dlfm`

</details>

<details><summary><b>Screenshots</b></summary>

#### Token mode:
Profile screenshot:

![Profile screenshot](https://i.imgur.com/SbtvFJa.png "Profile screenshot")

Screenshot from list of users:

![Screenshot from list of users](https://i.imgur.com/x5mWIXR.png "Screenshot from list of users") 

#### App mode:

Profile screenshot:

![Profile screenshot](https://i.imgur.com/shI5JJT.png "Profile screenshot")

Screenshot from list of users:

![Screenshot from list of users](https://i.imgur.com/F5E1GPz.png "Screenshot from list of users")
</details>

Config file (config.toml) have three parts: last.fm, discord and application settings.
Here are instruction to fill it correctly

## **1. Last.FM API Key**

Go to the [Last.FM API Page](https://www.last.fm/api/account/create) and sign in with your existing Last.FM username and password. It should bring you to the **Create API account** page and ask you for a few things.

It doesn't really matter what you put in most of the fields, but it should probably look something like this:

![LastFM Create API Account Screenshot](https://i.imgur.com/VQYa8nr.png?1)

After clicking Submit you should get a confirmation page with two items: *API Key* and *Shared Secret*. The API Key is the only one you need for this, but I recommend you save both for future use just in case, as they don't actually provide a way to retrieve these later.

![LastFM API Account Created Screenshot](https://i.imgur.com/1Qb7LeO.png "don't ask why names aren't same")

Copy and paste the API Key value into the config file in the `api_key = xxx` line **with quotes**

Then, fill field `username` with your last.fm username **with quotes**

Also you should fill field `check_interval` with integer number of seconds script should check new tracks. Recomended values: 5-15 seconds (if you are listening to grindcore, write 0). **Without quotes.**

## **2. Discord User Token**

Firstly, you should fill `use_app_mode` field with `true` or `false` (without quotes). By default it's `false` — token mode.

<details><summary>What is the difference?</summary>

| Token mode                  |               App mode |
|:----------------------------|-----------------------:|
| Using user token (unsafe)   |   Using application ID |
| Hard to get token           |         Easy to get ID |
| "Listening to ..."          |          "Playing ..." |
| Custom title                |         Title is fixed |
| No images                   | Large and small images |
| No image texts              |     Custom image texts |
| Can't work with text status | Works with text status |
</details>

Next step is different in App mode and Token mode, so I placed it in spoiler

<details><summary>Token mode</summary>

1. Go to Discord app
2. Press **Ctrl+R** (or **Cmd+R** on Mac)
3. Click the "*Application*" tab
4. Click and expand the "*Local Storage*" section
5. Click on the only entry in this section, *"https://discord.com/"*
6. Press **Ctrl+Shift+I** (or **Cmd+Shift+I** on Mac)
7. Wait few seconds
8. Right click -> Edit Value in the field to the right of "*token*"
9. Copy and paste the token value into the config file on the `token = discordtoken` line **saving quotes**.

![Desktop Token](https://i.imgur.com/sPs0New.png)
</details>

<details><summary>App mode</summary>

1. Go to [Discord Developer Portal's Applications page](https://discord.com/developers/applications "link")
2. Click to *"New application"* button
3. Type the name you want to see as title in your status (for example, type "qwe42" to see "Plaing qwe42" status)
4. Copy and paste *Application ID* into `app_id = 0123456789101112` line. **Delete quotes**!

![New application](https://i.imgur.com/Qd85IeE.png)

![Application ID](https://i.imgur.com/qphnFDa.png)
</details>

## **3. App settings**

If you want to ignore not critic errors, change "*endless_mode*" to `true`, else fill it with `false`. **Without quotes**!

Now you can modify Title in Token mode and Large image, Large text (what you'll see if you cover large image),
Small image and Small text (what you'll se if you cover small image) in App mode

In title, large image, large text, small text you can use these tags — parts of text that will be changed.
List of tags:

| Tag             |                     Value |    Type    |
|:----------------|---------------------------|-----------:|
| {{artist}}      | Track's artist(s) name(s) |    text    |
| {{album}}       | Track's album name        |    text    |
| {{name}}        | Track name                |    text    |
| {{album_image}} | Album cover from last.fm  |    image   |

For small image you can use other tags — they will be changed with icon of service:

<details><summary>Tags and images (too big to show without spoiler)</summary>

| Tag             |                                                                                                 Icon |
|:----------------|-----------------------------------------------------------------------------------------------------:|
| {{lastfm}}      | ![](http://icons.iconarchive.com/icons/danleech/simple/512/lastfm-icon.png)                          |
| {{deezer}}      | ![](https://www.macupdate.com/images/icons512/60905.png)                                             |
| {{youtube}}     | ![](https://seeklogo.com/images/Y/youtube-music-logo-50422973B2-seeklogo.com.png)                    |
| {{apple}}       | ![](http://ixd.prattsi.org/wp-content/uploads/2017/01/apple_music_logo_by_mattroxzworld-d982zrj.png) |
| {{vk}}          | ![](https://seeklogo.com/images/V/vk-icon-logo-10188561D5-seeklogo.com.png)                          |
| {{yandex}}      | ![](https://download.cdn.yandex.net/from/yandex.ru/support/ru/music/files/icon_main.png)             |
| {{soundcloud}}  | ![](https://icons.iconarchive.com/icons/sicons/basic-round-social/512/soundcloud-icon.png)           |

*icons aren't mine, i just tooked their URLs*
</details>


You can also fill `large_image` and `small_image` with your own URLs to any picture


## When you're done

Save your config file as "*config.toml*" and it should look something like this:

![Finished Config File](https://i.imgur.com/JVJOVy8.png)

Now just run the executable. It should connect to Discord and immediately start setting your status to whatever you're listening to on last.fm

If it's working, console will look like this:

![Running Executable](https://i.imgur.com/uDjruCs.png)