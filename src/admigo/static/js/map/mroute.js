;"use strict";
class Mroute {
  constructor(oin_in) {
    this.oin = oin_in;

    this.ros = {
      'el': document.getElementById('route-start'),
      'who': document.getElementById('ros-who'),
      'inp': document.getElementById('ros-inp'),
      'inp_hid': document.getElementById('ros-inp-lo'),
      'bma': document.getElementById('make-route'),
      'bcl': document.getElementById('clo-route'),
      'bde': document.getElementById('del-route'),
      'ma': null
    };

    this.ros.bcl.addEventListener('click', () => {
      this.close_route_form();
    });

    this.ros.bma.addEventListener('click', this.make_route_click.bind(this));
    this.ros.inp_hid.addEventListener('keyup', this.inp_hid_route_keyup.bind(this));
    this.ros.inp.addEventListener('input', this.inp_input.bind(this));
    this.ros.bde.addEventListener('click', this.del_route_click.bind(this));
  }

  cima(latlng) {
    let ret = L.circleMarker(latlng, {
      'stroke': true,
      'fill': true,
      'fillOpacity': 1,
      'fillColor': '#ffffff',
      'weight': 3,
      'color': '#2b5de5',
      'radius': 5
    }).addTo(this.oin.map);

    return ret;
  }

  map_route(geo_json) {
    let ro = L.geoJSON(geo_json, {
      'style': {
        'color': '#2b5de5',
        'weight': 10,
        'opacity': 0.5
      }
    });

    return ro;
  }

  show_route(some) {
    if (!some.ros) return;
    if (!some.ros.ma) return;
    if (!some.ros.ro) return;

    some.ros.ma.addTo(this.oin.map);
    some.ros.ro.addTo(this.oin.map);
  }

  move_start_route(some) {
    if (!some.ros) return;
    if (!some.ros.ma) return;

    this.oin.map.setView(some.ros.ma.getLatLng());
  }

  set_route_cid(some) {
    let nik = some.nik;
    this.ros.who.innerHTML = nik;
    this.ros.el.setAttribute('data-cid', some.cid);
    this.ros.inp.focus();

    this.show_route(some);
    this.move_start_route(some);
  }

  get_route_cid() {
    let cid = this.ros.el.getAttribute('data-cid');

    return cid;
  }

  close_route_form() {
    this.ros.el.setAttribute('data-cid', '');
    this.ros.who.innerHTML = '';

    if (this.ros.ma) {
      this.oin.map.removeLayer(this.ros.ma);
      this.ros.ma = null;
    }

    this.oin.sync_litems();
    this.oin.show_users();
  }

  clear_route(some) {
    if (!some.ros) {
      some.ros = {
        'ma': null,
        'ro': null,
        'ds': -1
      }
      return;
    }
    if (some.ros.ma) this.oin.map.removeLayer(some.ros.ma);
    if (some.ros.ro) this.oin.map.removeLayer(some.ros.ro);

    some.ros.ma = null;
    some.ros.ro = null;
    some.ros.ds = -1;
  }

  del_route_click() {
    const cid = this.get_route_cid();

    if (cid.length == 0) return;

    const some = this.oin.uslist[cid];

    if (!some) return;

    this.clear_route(some);
    this.close_route_form();
  }

  promise_route(st, en) {
    const url = this.oin.wsmap.ws.routeurl;
    const obj = {
      'coordinates': [[st.lng, st.lat], [en.lng, en.lat]]
    }
    const he = {
      'Content-Type': 'application/json; charset=utf-8',
      'Accept': 'application/json, application/geo+json, application/gpx+xml, img/png; charset=utf-8',
      'Authorization': this.oin.wsmap.ws.routekey
    }

    return axios.post(url, obj, {
      'headers': he
    });
  }

  get_dist(geo_json) {
    let sm = 0;

    geo_json.features.forEach((f) => {
      sm += parseFloat(f.properties.summary.distance);
    });

    return sm;
  }

  make_route(some, ma, cb_fun) {
    if (!some) return;
    if (!ma) return;
    if (!some.pos) return;

    const rm = ma.getLatLng();

    this.promise_route(some.pos, rm)
      .then((re) => {
        this.clear_route(some);

        some.ros.ma = ma;
        some.ros.ro = this.map_route(re.data);
        some.ros.ds = this.get_dist(re.data);

        this.close_route_form();
        this.show_route(some);
        if (cb_fun) cb_fun();
      })
      .catch((err) => {
        this.oin.showLog(err, true);
        if (cb_fun) cb_fun();
      });
  }

  make_route_click() {
    const cid = this.get_route_cid();

    if (cid.length == 0) return;

    const some = this.oin.uslist[cid];
    this.make_route(some, this.ros.ma, null);
  }

  bet_route(cids) {
    let some = this.oin.uslist[cids[0]];
    let euse = this.oin.uslist[cids[1]];
    let ma = this.cima(euse.ma.getLatLng());

    this.make_route(some, ma, this.oin.clear_bero.bind(this.oin));
  }

  inp_hid_route_keyup(ev) {
    const inp = ev.currentTarget;
    const val = inp.value;

    let ar = this.oin.str_to_latlng(val);
    if (!ar) return true;

    if (this.ros.ma) {
      this.oin.map.removeLayer(this.ros.ma);
      this.ros.ma = null;
    }

    const ln = L.latLng([ar[0], ar[1]]);
    let ma = this.cima(ln);

    this.ros.ma = ma;

    this.oin.map.setView(this.ros.ma.getLatLng());

    return true;
  }

  inp_input(ev) {
    const inp = ev.currentTarget;
    const val = inp.value;

    let ar = this.oin.str_to_latlng(val);
    if (!ar) return true;

    this.ros.inp_hid.value = val;
    this.oin.fun.trigger(this.ros.inp_hid, 'keyup');
  }

  map_click_route(e) {
    let cid = this.get_route_cid();

    if (cid.length == 0) return false;

    let elal = e.latlng
    this.ros.inp.value = `${elal.lat.toFixed(4)},${elal.lng.toFixed(4)}`;

    if (this.ros.ma) {
      this.oin.map.removeLayer(this.ros.ma);
      this.ros.ma = null;
    }

    this.ros.ma = this.cima(e.latlng);
  }

  init_map_event() {
    this.oin.map.on('click', this.map_click_route.bind(this));
  }
}
