$.ajax({
  type: "GET",
  url: "/buses",
  headers: {
    "Content-Type": "application/json",
  },
  dataType: "json",
  success: function (responseData, textStatus, jqXHR) {
    loadMap(responseData);
  },
  error: function (responseData, textStatus, errorThrown) {
    alert("request failed.");
  },
});

function loadMap(responseData) {
  var busData = parseResponse(responseData);
  var locateMyPositionControl = addLocateMyPositionControl();

  map = new ol.Map({
    target: "map",
    layers: [
      new ol.layer.Tile({
        source: new ol.source.OSM(),
      }),
      createBusLayer(busData),
    ],
    view: new ol.View(loadMapState()),
  });

  map.addControl(
    new ol.control.Control({
      element: locateMyPositionControl,
    })
  );

  subscribeAllEvents(map, locateMyPositionControl);
}

function parseResponse(responseData) {
  const features = [];
  for (let i = 0; i < responseData.length; i++) {
    var bus = responseData[i];
    features.push(
      new ol.Feature({
        geometry: new ol.geom.Point(
          ol.proj.fromLonLat([bus.longitude, bus.latitude])
        ),
        text: bus.bus_service_name,
        busRef: bus.vehicleRef,
        inService: bus.inService,
        details: `Direction: ${bus.originName} &rarr; ${bus.destinationName} </br> Aimed time: ${bus.originAimedDepartureTime} &rarr; ${bus.destinationAimedArrivalTime}`,
      })
    );
  }
  return features;
}

function createBusLayer(features) {
  return new ol.layer.Vector({
    source: new ol.source.Vector({
      features,
    }),
    style: function (feature) {
      return new ol.style.Style({
        image: new ol.style.Circle({
          radius: 12,
          stroke: new ol.style.Stroke({
            color: "black",
          }),
          fill: new ol.style.Fill({
            color: feature.get("inService") === true ? "#3be62c" : "#e35922",
          }),
        }),
        text: new ol.style.Text({
          text: feature.get("text"),
          font: "10px Calibri,sans-serif",
          fill: new ol.style.Fill({ color: "#000" }),
          offsetY: 0,
        }),
      });
    },
  });
}

function addLocateMyPositionControl() {
  const locate = document.createElement("div");
  locate.className = "ol-control ol-unselectable locate";
  locate.innerHTML = '<button title="Locate my position">&#9673</button>';
  return locate;
}

function loadMapState() {
  const savedCenter = localStorage.getItem("mapCenter");
  const savedZoom = localStorage.getItem("mapZoom");
  return {
    center: savedCenter
      ? ol.proj.fromLonLat(JSON.parse(savedCenter))
      : ol.proj.fromLonLat([-1.266796088748171, 51.6387863304397]),
    zoom: savedZoom ? parseFloat(savedZoom) : 12.5,
  };
}

function subscribeAllEvents(map, locateControl) {
  map.on("click", function (event) {
    map.forEachFeatureAtPixel(event.pixel, function (feature) {
      serviceMessage =
        feature.get("inService") === true
          ? "in service"
          : "&#9888; inactive or not ready to depart";

      $("#toast-body-title").text(
        `${feature.get("text")} (Ref: ${feature.get("busRef")})`
      );
      $("#toast-body-content").html(feature.get("details"));
      $("#service_info").html(serviceMessage);

      bootstrap.Toast.getOrCreateInstance(
        document.getElementById("liveToast")
      ).show();
    });
  });

  map.getView().on("change", function () {
    const view = map.getView();
    const center = ol.proj.toLonLat(view.getCenter());
    const zoom = view.getZoom();

    localStorage.setItem("mapCenter", JSON.stringify(center));
    localStorage.setItem("mapZoom", zoom);
  });

  locateControl.addEventListener("click", function () {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const userCoordinates = ol.proj.fromLonLat([
            position.coords.longitude,
            position.coords.latitude,
          ]);

          map.getView().setCenter(userCoordinates);
          map.getView().setZoom(15);
        },
        (error) => {
          alert("Error getting location: " + error.message);
        }
      );
    } else {
      alert("Geolocation is not supported by your browser.");
    }
  });
}
