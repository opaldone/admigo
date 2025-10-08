;"use strict";
class Amap {
  constructor() {
    this.map = null;
    this.wsmap = null;
    this.taber = null;
    this.ulo = null;
    this.elmap = document.getElementById('map');
    this.uslist = {};

    this.lg = document.getElementById('ws-logs');
    this.lg_clear = document.getElementById('ws-clear');
    this.lg_errors = document.getElementById('errors-cnt');

    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));

    this.lg_clear.addEventListener('click', this.lg_clear_click.bind(this));

    this.ros = {
      'el': document.getElementById('route-start'),
      'who': document.getElementById('ros-who'),
      'inp': document.getElementById('ros-inp'),
      'bma': document.getElementById('make-route'),
      'bcl': document.getElementById('clo-route'),
      'bde': document.getElementById('del-route'),
      'ma': null
    };

    this.ros.bcl.addEventListener('click', () => {
      this.close_route_form();
    });

    this.ros.bma.addEventListener('click', this.make_route_click.bind(this));
    this.ros.inp.addEventListener('input', this.inp_route_update.bind(this));
    this.ros.bde.addEventListener('click', this.del_route_click.bind(this));
  }

  fm_tm(co) {
    return co < 10 ? '0' + co : co;
  }

  get_tm() {
    let nw = new Date();
    let ho = nw.getHours();
    let mi = nw.getMinutes();
    let se = nw.getSeconds();

    return this.fm_tm(ho) + ':' + this.fm_tm(mi) + ':' + this.fm_tm(se);
  }

  ref_log_cnt() {
    this.lg_errors.innerHTML = '';
    let errs = this.lg.querySelectorAll('.err');
    let cc = errs.length;

    if (cc == 0) return;

    this.lg_errors.textContent = cc;
  }

  lg_clear_click(e) {
    e.preventDefault();
    e.stopPropagation();

    this.lg.innerHTML = '';
    this.ref_log_cnt();

    return false;
  }

  showLog(msg, err) {
    let si = '<li';
    if (err) {
      si += ' class="err"';
    }
    si += '><span class="lg-msg">' + msg + '</span>' +
      '<span class="lg-tm">' + this.get_tm() + '</span></li>';

    let tem = document.createElement('template');
    tem.innerHTML = si;

    this.lg.prepend(tem.content);

    setTimeout(() => {
      this.ref_log_cnt();
    }, 100);

    return false;
  }

  show_route(some) {
    if (!some.ros) return;
    if (!some.ros.ma) return;
    if (!some.ros.ro) return;

    some.ros.ma.addTo(this.map);
    some.ros.ro.addTo(this.map);
  }

  move_start_route(some) {
    if (!some.ros) return;
    if (!some.ros.ma) return;

    this.map.setView(some.ros.ma.getLatLng(), 17);
  }

  set_route_cid(cid) {
    const some = this.uslist[cid];
    let nik = some.nik;
    this.ros.who.innerHTML = nik;
    this.ros.el.setAttribute('data-cid', cid);

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
    this.ros.inp.value = '';
    if (this.ros.ma) {
      this.map.removeLayer(this.ros.ma);
      this.ros.ma = null;
    }
    this.ulo.sync_litems();
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
    if (!some.ros.ma) return;
    if (!some.ros.ro) return;

    this.map.removeLayer(some.ros.ma);
    this.map.removeLayer(some.ros.ro);

    some.ros.ma = null;
    some.ros.ro = null;
    some.ros.ds = -1;
  }

  del_route_click() {
    const cid = this.get_route_cid();

    if (cid.length == 0) return;

    const some = this.uslist[cid];

    if (!some) return;

    this.clear_route(some);
    this.close_route_form();
  }

  make_route_click() {
    const cid = this.get_route_cid();

    if (cid.length == 0) return;
    if (!this.ros.ma) return;

    const some = this.uslist[cid];

    if (!some) return;

    const pos = some.pos;

    if (!pos) return;

    const rm = this.ros.ma.getLatLng();

    const url = 'https://api.openrouteservice.org/v2/directions/driving-car/geojson';
    const obj = {
      'coordinates': [[rm.lng, rm.lat], [pos.lng, pos.lat]]
    }
    const he = {
      'Content-Type': 'application/json; charset=utf-8',
      'Accept': 'application/json, application/geo+json, application/gpx+xml, img/png; charset=utf-8',
      'Authorization': 'eyJvcmciOiI1YjNjZTM1OTc4NTExMTAwMDFjZjYyNDgiLCJpZCI6IjI5Y2U4YjE0YmMwZTQ2ZDVhNDI3NzFlNDU2MzhlODI5IiwiaCI6Im11cm11cjY0In0='
    }

    axios.post(url, obj, {
      'headers': he
    })
      .then((re) => {
        this.clear_route(some);

        some.ros.ma = this.ros.ma;
        some.ros.ro = L.geoJSON(re.data);

        let sm = 0;
        re.data.features.forEach((f) => {
          sm += parseFloat(f.properties.summary.distance);
        });

        some.ros.ds = sm;

        this.close_route_form();
        this.show_route(some);
      })
      .catch((err) => {
        console.log(err);
      });
  }

  cima(latlng) {
    let ret = L.circleMarker(latlng, {
      'stroke': false,
      'fill': true,
      'fillOpacity': 1,
      'fillColor': '#2C5DE5',
      'radius': 5
    }).addTo(this.map);

    return ret;
  }

  inp_route_update(ev) {
    const inp = ev.currentTarget;
    const val = inp.value;

    let ar = val.split(/[,;: ]/);

    console.log('inp_route_update', val, ar);

    if (ar.length != 2) {
      inp.value = '';
      return;
    }

    if (this.ros.ma) {
      this.map.removeLayer(this.ros.ma);
      this.ros.ma = null;
    }

    const ln = L.latLng([ar[0], ar[1]]);
    let ma = this.cima(ln);

    this.ros.ma = ma;

    this.map.setView(this.ros.ma.getLatLng(), 17);
  }

  map_click(e) {
    let cid = this.get_route_cid();

    if (cid.length == 0) return false;

    let elal = e.latlng
    this.ros.inp.value = `${elal.lat.toFixed(4)},${elal.lng.toFixed(4)}`;

    if (this.ros.ma) {
      this.map.removeLayer(this.ros.ma);
      this.ros.ma = null;
    }

    this.ros.ma = this.cima(e.latlng);
  }

  set_wsmap() {
    this.taber = new Taber();
    this.wsmap = new Wsmap(this);
    this.ulo = new Uloca(this, this.fun);
  }

  init_map() {
    this.map = L.map(this.elmap).setView([57.989287, 56.213889], 13);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(this.map);

    this.map.on('click', this.map_click.bind(this));

    this.set_wsmap();
  }

  handler() {
    this.init_map();
  }

  set_uslist_item(v) {
    let cid = v.cid;

    if (!this.uslist[cid]) {
      this.uslist[cid] = {
        'nik': '',
        'issender': false,
        'pos': null
      };
    }

    this.uslist[cid]['nik'] = v.nik;
    this.uslist[cid]['issender'] = v.issender;
    this.uslist[cid].pos = v.pos;
  }

  rem_uslist_item(v) {
    const cid = v.cid;

    this.ulo.rem_from_list(cid);

    if (!this.uslist[cid]) return;

    if (this.uslist[cid].ma) {
      this.map.removeLayer(this.uslist[cid].ma);
    }

    if (this.uslist[cid].ci) {
      this.map.removeLayer(this.uslist[cid].ci);
    }

    delete this.uslist[cid];
  }

  ref_uslist(cont) {
    let js = JSON.parse(cont);

    if (!js) {
      return this.showLog("Failed to parse ref_uslist", true)
    }

    js.forEach((v, _) => {
      if (!v.issender) return;

      this.set_uslist_item(v)
    });

    this.ulo.ref_list();
  }

  sender_hi(cont) {
    let js = JSON.parse(cont);

    if (!js) {
      return this.showLog("Failed to parse sender_hi", true);
    }

    this.set_uslist_item(js);
    this.ulo.ref_list();
  }

  sender_stop(cont) {
    let js = JSON.parse(cont);

    if (!js) {
      return this.showLog("Failed to parse sender_stop", true);
    }

    this.rem_uslist_item(js);
  }

  req_loca(cid) {
    this.wsmap.req_loca_cid(cid);
  }

  ans_loca(cont) {
    let js = JSON.parse(cont);

    if (!js) {
      return this.showLog("Failed to parse ans_loca", true);
    }

    this.set_uslist_item(js);
    this.ulo.ref_ma(js.cid);
  }
}

new Amap();
