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
        details: `Direction: ${bus.originName} &rarr; ${bus.destinationName} </br> Aimed time: ${bus.originAimedDepartureTime} &rarr; ${bus.destinationAimedArrivalTime}`
      })
    );
  }

  const busLayer = new ol.layer.Vector({
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

  const locate = document.createElement("div");
  locate.className = "ol-control ol-unselectable locate";
  locate.innerHTML = '<button title="Locate my position">&#9673</button>';
  locate.addEventListener("click", function () {
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

  function loadMapState() {
    const savedCenter = localStorage.getItem("mapCenter");
    const savedZoom = localStorage.getItem("mapZoom");

    if (savedCenter && savedZoom) {
      return {
        center: JSON.parse(savedCenter),
        zoom: parseFloat(savedZoom),
      };
    }
    return null;
  }

  const mapDefaultSettings = {
    center: [-1.271189, 51.605066],
    zoom: 15,
  };

  const savedState = loadMapState();
  const mapCenter = savedState
    ? ol.proj.fromLonLat(savedState.center)
    : ol.proj.fromLonLat(mapDefaultSettings.center);
  const initialZoom = savedState ? savedState.zoom : mapDefaultSettings.zoom;

  const map = new ol.Map({
    target: "map",
    layers: [
      new ol.layer.Tile({
        source: new ol.source.OSM(),
      }),
      busLayer,
    ],
    view: new ol.View({
      center: mapCenter,
      zoom: initialZoom,
    }),
  });

  map.addControl(
    new ol.control.Control({
      element: locate,
    })
  );

  map.on("click", function (event) {
    map.forEachFeatureAtPixel(event.pixel, function (feature) {
      $("#busInfoTitle").text(`Bus info: ${feature.get("text")} (Ref: ${feature.get("busRef")})`);
      $("#model_body_main_content").html(feature.get("details"));
      if (feature.get("inService") === false) {
          $("#model_body_warning_content").html("&#9888; Arrival time is over. Bus might be not in service BUT waiting to depart");
      }

      $("#busInfoModel").modal("show");
    });
  });

  map.getView().on("change", function () {
    const view = map.getView();
    const center = ol.proj.toLonLat(view.getCenter());
    const zoom = view.getZoom();

    localStorage.setItem("mapCenter", JSON.stringify(center));
    localStorage.setItem("mapZoom", zoom);
  });
}
