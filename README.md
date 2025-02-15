# UK Realtime Bus

["UK Realtime bus"](https://link-url-here.org) is a simple web appication which provides realtime bus information in UK. 
The source of data is public available API managed by Gov.UK bus data: https://data.bus-data.dft.gov.uk

The "UK Realtime Bus" application displays this data on map so users can observe available buses. Google maps doesn't provide such functions which is sometimes inconvinient when a bus is expected to come on time but doesn't appear.

Right now the functionality is pretty limited and aimed display buses only for a single provider which is "Oxford Bus Company and Thames Travel" because the main motivation was to manage personal schedule issues. Later, this might be improved to cover whole UK bus providers.

Feel free to contact kutuev93@yandex.ru if you have any ideas or suggetions

# Deployment status

Currently the application is served on Google Cloud and already available here: https://uk-realtime-bus-858017075201.us-central1.run.app


# Changelog

This is lisf of added/planned functional features:

- [x] Added button to centerlise map based on user position 
- [x] Automatically persist user map position and zoom level in localStorage
- 
