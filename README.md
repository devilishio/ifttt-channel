# IFTTT Webhook Channel
This is a small proxy service for invoking a webhook via an IFTTT recipe.

It works by emulating the part of the Wordpress XMLRPC API so that you can add a Wordpress channel to IFTTT that points at the service. You then specify the webhook to call in the recipe content. It is implemented in Go.

_Note:_ this is a work in progress and is not yet even remotely functional. Please look at https://github.com/femto113/node-ifttt-webhook if you want something that can do the same now.
