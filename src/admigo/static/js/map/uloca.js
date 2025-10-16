;"use strict";
class Uloca {
  constructor(oin_in, fun_in) {
    this.oin = oin_in;
    this.fun = fun_in;
    this.list = document.getElementById('map-users');
    this.user_cnt = document.getElementById('users-cnt');
    this.li_tag = this.get_li_tag();
    this.bero = [];
  }

  docon() {
    this.list.querySelectorAll('li').forEach(el => {
      if (this.oin.fun.once(el, 'li_click')) return;
      el.addEventListener('click', this.li_click.bind(this));
    });

    this.list.querySelectorAll('.us-route').forEach(el => {
      if (this.oin.fun.once(el, 'us_route_click')) return;
      el.addEventListener('click', this.us_route_click.bind(this));
    });

    this.list.querySelectorAll('.bet-us-route').forEach(el => {
      if (this.oin.fun.once(el, 'bet_us_route_click')) return;
      el.addEventListener('click', this.bet_us_route_click.bind(this));
    });
  }

  ref_coo_cont(some) {
    let liel = document.getElementById(some.cid);

    if (!liel) return;

    let sid = 'coo-' + some.cid;
    let el = document.getElementById(sid);

    if (!el) return;

    let msg = 'Press the item to get the location';

    if (some.pos) {
      msg = `<i class="fa-solid fa-location-crosshairs"></i><div class="i-val">${some.pos.lat.toFixed(4)},${some.pos.lng.toFixed(4)}</div>`;
    }

    el.innerHTML = msg;
  }

  fm_distance(di) {
    if (di <= 0) return '';

    const kms = Math.floor(di / 1000);
    const mts = di % 1000;

    let ret = "";

    if (kms > 0) {
      ret += `${kms}<span>km</span>`;
    }

    if (mts > 0) {
      ret += ` ${mts.toFixed(0)}<span>m</span>`;
    }

    if (kms === 0 && mts === 0) {
      ret = "";
    }

    return ret;
  }

  ref_dista_cont(some) {
    if (!some.ros) return;
    if (!some.ros.ds) return;

    let liel = document.getElementById(some.cid);
    if (!liel) return;
    let sid = 'dista-' + some.cid;
    let el = document.getElementById(sid);
    if (!el) return;

    const dis = this.fm_distance(some.ros.ds);

    let msg = '';
    if (dis.length > 0) {
      msg = '<i class="fa-solid fa-compass-drafting"></i><div class="i-val">' + dis + '</div>';
    }

    el.innerHTML = msg;
  }

  sync_litems() {
    const cids = Object.keys(this.oin.uslist);

    if (cids.length == 0) return;

    const rcid = this.oin.get_route_cid();

    cids.forEach((cid, _) => {
      let litem = document.getElementById(cid);
      const some = this.oin.uslist[cid];

      if (some.in_se) {
        litem.classList.add('in-se');
      } else {
        litem.classList.remove('in-se');
      }

      if (cid == rcid) {
        litem.classList.add('in-route');
      } else {
        litem.classList.remove('in-route');
      }

      if (some.ros && some.ros.ro) {
        litem.classList.add('w-route');
      } else {
        litem.classList.remove('w-route');
      }

      if (this.bero.includes(cid)) {
        litem.classList.add('in-bero')
      } else {
        litem.classList.remove('in-bero');
      }

      this.ref_dista_cont(some);
    });
  }

  clear_timer(some) {
    if (!some.tm) return;

    clearTimeout(some.tm);
    some.tm = null;
  }

  update_timer(some) {
    this.clear_timer(some);

    if (!some.in_se) return;

    some.tm = setTimeout(() => {
      this.oin.req_loca(some);
    }, 5000);
  }

  ref_ma(cid) {
    let some = this.oin.uslist[cid];

    if (!some) return false;

    this.ref_coo_cont(some);
    this.update_timer(some);

    if (!some.pos) return false;

    let sp = [some.pos.lat, some.pos.lng];
    let can_move = this.oin.get_route_cid().length == 0;

    if (some.ma) {
      some.ma.setLatLng(sp);
      if (can_move) {
        some.ma.openPopup();
      } else {
        some.ma.closePopup();
      }
    } else {
      some.ma = L.marker(sp, {
        'icon': this.oin.mai
      }).addTo(this.oin.map);
      let pop = some.nik ? some.nik : some.cid;
      some.ma.bindPopup(pop).openPopup();
    }

    if (some.ci) {
      some.ci.setLatLng(sp);
      some.ci.setRadius(some.pos.acc);
    } else {
      some.ci = L.circle(sp, {
        'fillColor': '#2b5de5',
        'fillOpacity': 0.2,
        'stroke': false,
        'color': '#2b5de5',
        'weight': 1,
        'radius': some.pos.acc
      }).addTo(this.oin.map);
    }

    if (can_move) {
      this.oin.move_to_ma(some);
    }
  }

  clear_in_se() {
    const cids = Object.keys(this.oin.uslist);

    if (cids.length == 0) return;

    const cid = cids.find((cc) => {
      return this.oin.uslist[cc].in_se == true;
    });

    if (!cid) return;

    this.oin.uslist[cid].in_se = false;
  }

  li_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget;
    let cid = el.getAttribute('id');
    let some = this.oin.uslist[cid];

    if (!some) return;

    if (some.in_se) {
      this.clear_timer(some);
      some.in_se = false;
      this.sync_litems();
      return;
    }

    this.clear_in_se();
    some.in_se = true;
    this.sync_litems();

    this.oin.move_to_ma(some);
    this.oin.req_loca(some);

    return false;
  }

  us_route_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget;
    let li_cont = this.fun.parent(el, '.map-us-li');
    let cid = li_cont.getAttribute('id');
    let in_route = li_cont.classList.contains('in-route');
    let some = this.oin.uslist[cid];

    if (!some.pos) {
      this.ref_coo_cont(some);
      return false;
    }

    if (in_route) {
      this.oin.close_route_form();
      return false;
    }

    this.oin.set_route_cid(some);
    this.sync_litems();
    this.oin.hide_tabs();

    return false;
  }

  check_bero() {
    if (this.bero.length < 2) return;

    this.oin.bet_route(this.bero);
  }

  clear_bero() {
    this.bero = [];
    this.sync_litems();
  }

  bet_us_route_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget;
    let li_cont = this.fun.parent(el, '.map-us-li');
    let cid = li_cont.getAttribute('id');
    let in_bero = li_cont.classList.contains('in-bero');
    let some = this.oin.uslist[cid];

    if (!some.ma) {
      this.ref_coo_cont(some);
      return false;
    }

    if (in_bero) {
      this.bero = this.bero.filter(cidin => cidin !== cid);
      this.sync_litems();
      return false;
    }

    if (!this.bero.includes(cid)) {
      this.bero.push(cid);
    }

    this.check_bero();
    this.sync_litems();

    return false;
  }

  ref_cnt() {
    this.user_cnt.innerHTML = '';
    let cc = this.list.children.length;

    if (cc == 0) return;

    this.user_cnt.textContent = cc;
  }

  get_li_tag() {
    let ret = '<li id="#CID#" class="map-us-li">' +
      '<div class="map-us-cont">' +
      '<div id="nik-#CID#" class="nik-cont"></div>' +
      '<div class="info-cont">' +
      '<div class="coo-cont" id="coo-#CID#" title="Location"></div>' +
      '<div class="coo-cont" id="dista-#CID#" title="Distance"></div>' +
      '</div>' +
      '<div class="us-btn-cont">' +
      '<span class="us-route" title="Make a route">' +
      '<i class="fa-solid fa-car-side"></i>' +
      '</span>' +
      '<span class="bet-us-route" title="Make a route between people">' +
      '<i class="fa-solid fa-person-walking-arrow-loop-left"></i>' +
      '</span>' +
      '</div>' +
      '</div>' +
      '</li>';

    return ret;
  }

  get_new_li(cid) {
    let si = this.li_tag
      .replace(/#CID#/g, cid);

    let tem = document.createElement('template');
    tem.innerHTML = si;

    return tem.content;
  }

  ref_list() {
    const cids = Object.keys(this.oin.uslist);

    if (cids.length == 0) return;

    cids.forEach((cid, _) => {
      let nel = document.getElementById(cid);
      const some = this.oin.uslist[cid];

      if (!nel) {
        nel = this.get_new_li(cid);
        this.list.appendChild(nel);
      }

      let nik = document.getElementById('nik-' + cid);
      nik.textContent = some.nik.length ? some.nik : some.cid;
    });

    this.docon();

    this.ref_cnt();
  }

  rem_from_list(cid) {
    let nel = document.getElementById(cid);

    if (!nel) return;

    nel.remove();

    this.ref_cnt();
  }
}
