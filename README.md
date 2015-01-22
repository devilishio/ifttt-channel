# IFTTT Webhook Channel
This is a small proxy service for invoking a webhook via an IFTTT recipe.

It works by emulating the part of the Wordpress XMLRPC API so that you can add a Wordpress channel to IFTTT that points at the service. You then specify the webhook to call in the recipe content. It is implemented in Go.

_Note:_ This is a work in progress. At the moment it does very basic proxying using a URL that you put in the Wordpress description field. It currently does this as a GET request. Later it will support POSTs and authentication and probably will not allow proxying to any URL, just configured base URLs. So basically in it's current state it's good for playing around and learning but don't use it for anything serious and don't leave it exposed on the internet for any length of time! 
