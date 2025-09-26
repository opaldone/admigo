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

    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));

    this.lg_clear.addEventListener('click', this.ws_clear_click.bind(this));
  }

  ws_clear_click(e) {
    e.preventDefault();
    e.stopPropagation();

    this.lg.innerHTML = '';

    return false;
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

    return false;
  }

  set_wsmap() {
    this.taber = new Taber();
    this.wsmap = new Wsmap(this);
    this.ulo = new Uloca(this);
  }

  init_map() {
    this.map = L.map(this.elmap).setView([57.989287, 56.213889], 13);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(this.map);

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
