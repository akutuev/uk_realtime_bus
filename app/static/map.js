$.ajax({
    type: 'GET',
    url: '/buses',
    headers: {
        'Content-Type': 'application/json',
    },
    dataType: 'json',
    success: function(responseData, textStatus, jqXHR) {
        loadMap(responseData);
    },
    error: function (responseData, textStatus, errorThrown) {
        alert('request failed.');
    }
});

function loadMap(data) {
    console.dir(data);   

    const features = [];
    for(let i = 0; i < data.length; i++) {
        var bus = data[i];

        var warning = bus.inService === true ? "" : "<br/><center> &#9888<b>Arrival time is over. Bus might be not in service BUT waiting to depart</b> <center>"
        
        features.push(
            new ol.Feature({
                geometry: new ol.geom.Point(
                    ol.proj.fromLonLat([
                    bus.longitude, bus.latitude
                    ])),
                text: bus.bus_service_name,
                busRef: bus.vehicleRef,
                inService: bus.inService,
                details: 
                    "Direction: " + bus.originName + " &rarr; " + bus.destinationName + "</br>" +
                    "Aimed time: " + bus.originAimedDepartureTime + " &rarr; " + bus.destinationAimedArrivalTime + "</br>" +
                    warning        
            }));
        }

    const vectorSource = new ol.source.Vector({
        features
    });

    const vectorLayer = new ol.layer.Vector({
        source: vectorSource,
        style: function (feature) {
            return new ol.style.Style({
                image: new ol.style.Circle({
                    radius: 12,
                    stroke: new ol.style.Stroke({
                        color: 'black'
                    }),
                    fill: new ol.style.Fill({
                        color: feature.get('inService') === true ? '#3be62c' : "#e35922"
                    })
                }),
                text: new ol.style.Text({
                        text: feature.get('text'),
                        font: '10px Calibri,sans-serif',
                        fill: new ol.style.Fill({ color: '#000' }),
                        offsetY: 0,
                    })
                });
        }
    });
    

    const locate = document.createElement('div');
    locate.className = 'ol-control ol-unselectable locate';
    locate.innerHTML = '<button title="Locate my position">&#9673</button>';
    locate.addEventListener('click', function () {
        locateUser();
      });

    const map = new ol.Map({
        target: 'map',
        layers: [
            new ol.layer.Tile({
                source: new ol.source.OSM()
            }),
            vectorLayer
        ],
        view: new ol.View({
            center: ol.proj.fromLonLat([-1.271189,51.605066]),
            zoom: 15
        })
    });

    map.addControl(
        new ol.control.Control({
          element: locate,
        }),
      );

    map.on('click', function(event) {
        map.forEachFeatureAtPixel(event.pixel, function(feature,layer) {
            document.getElementById("busInfoModelLabel").innerHTML = "Bus info: " + feature.get('text') + " (Ref: " + feature.get('busRef') + ")";
            document.getElementById("modal_body").innerHTML = feature.get('details');
            $('#busInfoModel').modal('show');
        });
    }); 

    function locateUser() {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition((position) => {
                const userCoordinates = ol.proj.fromLonLat([
                    position.coords.longitude,
                    position.coords.latitude
                ]);

                map.getView().setCenter(userCoordinates);
                map.getView().setZoom(15);
            }, (error) => {
                alert('Error getting location: ' + error.message);
            });
        } else {
            alert('Geolocation is not supported by your browser.');
        }
    }
}
