# fontface-lab-server
Server for the fontface-lab browser extension to query for all available fonts in the Google Fonts API. 
This server has two main purposes:

1. Serves as an alternative to the https://fonts.google.com/metadata/fonts endpoint, since that endpoint enforces the same-origin policy.
2. Data wrangling of font data ahead of time for the client.
   
# Overview

The API boilerplate lives in `main.go`, and code specific to fetching/wrangling the Google Fonts API lives in the `data` folder.

The flow of data is relatively simple:

1. Client hits `api/font-family-list` endpoint.
2. Server computes the list of fonts available in the Google Fonts API:
    1. If a cached result exists, return it and jump to step 3.
    2. Otherwise:
        1. Hit the https://fonts.google.com/metadata/fonts endpoint.
        2. Do some data wrangling of the data we get back.
        3. Store the data in the cache, return it, and jump to step 3.
3. Return to the client with the appropriate caching headers.
