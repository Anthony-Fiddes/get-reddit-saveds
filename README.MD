# get-reddit-saveds

This is a very simple script type application that gets a user's latest saved
posts from Reddit. By default, it tries to save 1,000 posts to output.json in
the same directory as the executable.

This is an okay example of how to start messing around with the Reddit api 
using Go!

## How to Use

This depends on Reddit's API. As such, you need to create a Reddit app. You can
find the instructions on how to do that
[here](https://github.com/reddit-archive/reddit/wiki/OAuth2-Quick-Start-Example).

Once you have all that information, you really only have to build the executable
and run it. You can accomplish that by running `go run .` in the same directory
as the main.go file.

## Options

You can change the constants in the source code to accomplish different
behaviors. `limit` controls how many posts are requested, and `settings`
and `output` can change where those respective files are looked for/created.

You can put a settings.json file in your directory so that you don't have to
specify your credentials every time you want to use this script. I created this
part of the script mainly to make my own development easier. I really don't
recommend that you store your passwords/secrets in plain text on your computer;
it's just a convenience feature.
