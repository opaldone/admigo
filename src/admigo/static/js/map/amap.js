;"use strict";
class Amap {
  constructor() {
    this.map = null;
    this.mai = null;
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
  }

  test_fill_uslist() {
    let uu = `{
        "cid": "111",
        "nik": "Шавкунов Павел Сергеевич",
        "issender": true
    }`;

    let lo = `{
      "cid": "111",
      "nik": "Шавкунов Павел Сергеевич",
      "pos": {
        "lat": 57.9895,
        "lng": 56.2143,
        "acc": 15
      }
    }`

    let uul = `{
        "cid": "222",
        "nik": "Some user 222",
        "issender": true
    }`;

    let lol = `{
      "cid": "222",
      "nik": "Some user 222",
      "pos": {
        "lat": 58.0202,
        "lng": 56.3009,
        "acc": 15
      }
    }`

    let uuq = `{
        "cid": "333",
        "nik": "Some user 333",
        "issender": true
    }`;

    let loq = `{
      "cid": "333",
      "nik": "Some user 333",
      "pos": {
        "lat": 57.9875,
        "lng": 56.2152,
        "acc": 30
      }
    }`

    setTimeout(() => {
      this.sender_hi(uu);
      setTimeout(() => {
        this.ans_loca(lo);
      }, 500)
    }, 1000);

    setTimeout(() => {
      this.sender_hi(uuq);
      setTimeout(() => {
        this.ans_loca(loq);
      }, 500)
    }, 2000);

    // setTimeout(() => {
      // this.sender_hi(uul);
      // setTimeout(() => {
        // this.ans_loca(lol);
      // }, 500)
    // }, 3000);
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

  fm_tm(co) {
    return co < 10 ? '0' + co : co;
  }

  str_to_latlng(str) {
    return str.split(/[,;: ]/);
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

  sync_litems() {
    this.ulo.sync_litems()
  }

  move_to_ma(some) {
    if (!some.ma) return;
    this.map.setView(some.ma.getLatLng());
    some.ma.openPopup();
  }

  set_wsmap() {
    this.taber = new Taber();
    this.wsmap = new Wsmap(this);
    this.ulo = new Uloca(this, this.fun);
    this.mro = new Mroute(this);
  }

  init_map() {
    this.map = L.map(this.elmap);
    this.map.on('load', () => {
      L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 18,
        attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
      }).addTo(this.map);

      this.wsmap.startWs();

      // this.test_fill_uslist();
    });
    this.mai = L.icon({
      iconUrl: '/static/images/map/ma.png',
      iconAnchor:   [10, 25],
      popupAnchor:  [0, -25]
    });
    this.mro.init_map_event();
    this.map.setView(this.str_to_latlng(this.wsmap.ws.startpoint), 17);
  }

  handler() {
    this.set_wsmap();
  }

  hide_tabs() {
    if (!this.taber) return;

    this.taber.hide_tabs();
  }

  set_uslist_item(v) {
    const cid = v.cid;

    if (!this.uslist[cid]) {
      this.uslist[cid] = {
        'cid': cid,
        'nik': '',
        'issender': false,
        'pos': null
      };
    }

    this.uslist[cid]['nik'] = v.nik;
    this.uslist[cid]['issender'] = v.issender;
    this.uslist[cid]['pos'] = v.pos;
  }

  rem_uslist_item(v) {
    const cid = v.cid;

    this.ulo.rem_from_list(cid);

    let some = this.uslist[cid]

    if (!some) return;

    if (some.ma) this.map.removeLayer(some.ma);
    if (some.ci) this.map.removeLayer(some.ci);
    this.mro.clear_route(some);

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

  req_loca(some) {
    this.wsmap.req_loca_cid(some.cid);
  }

  ans_loca(cont) {
    let js = JSON.parse(cont);

    if (!js) {
      return this.showLog("Failed to parse ans_loca", true);
    }

    this.set_uslist_item(js);
    this.ulo.ref_ma(js.cid);
  }

  set_route_cid(some) {
    this.mro.set_route_cid(some);
  }

  get_route_cid() {
    return this.mro.get_route_cid();
  }

  close_route_form() {
    this.mro.close_route_form();
  }

  bet_route(cids) {
    this.mro.bet_route(cids);
  }

  clear_bero() {
    this.ulo.clear_bero();
  }
}

new Amap();
