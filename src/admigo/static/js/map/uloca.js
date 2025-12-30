;"use strict";
class Uloca {
  constructor(oin_in, fun_in) {
    this.oin = oin_in;
    this.fun = fun_in;
    this.list = document.getElementById('map-users');
    this.user_cnt = document.getElementById('users-cnt');
    this.li_tag = this.get_li_tag();
    this.bero = [];
    this.time_out = 5000;
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

    this.list.querySelectorAll('.us-monit').forEach(el => {
      if (this.oin.fun.once(el, 'us_monit_click')) return;
      el.addEventListener('click', this.us_monit_click.bind(this));
    });

    this.list.querySelectorAll('.us-chat').forEach(el => {
      if (this.oin.fun.once(el, 'us_chat_click')) return;
      el.addEventListener('click', this.us_chat_click.bind(this));
    });
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

  fm_bat(ba) {
    let num = Math.ceil(ba / 10) * 10;

    switch (true) {
      case (num == 100):
        return 'full';
      case (num >= 75 && num < 100):
        return 'three-quarters';
      case (num >= 50 && num < 75):
        return 'half';
      case (num >= 25 && num < 50):
        return 'quarter';
      default:
        return 'empty';
    }
  }

  ref_info(elin, msg) {
    if (msg.length <= 0) return;
    elin.innerHTML = msg;
    elin.classList.add('sh');
  }

  ref_coo_cont(some) {
    let sid = 'coo-' + some.cid;
    let el = document.getElementById(sid);

    if (!el) return;

    let msg = '&mdash;';

    if (some.pos) {
      msg = `${some.pos.lat.toFixed(4)},${some.pos.lng.toFixed(4)}`;
    }

    msg = `<i class="fa-solid fa-location-crosshairs"></i><div class="i-val">${msg}</div>`;

    this.ref_info(el, msg);
  }

  ref_bat_cont(some) {
    if (!some.bat) return;

    let ba = parseInt(some.bat);
    if (ba < 0) return;

    let sid = 'bat-' + some.cid;
    let el = document.getElementById(sid);
    if (!el) return;

    let bc = this.fm_bat(ba);

    let msg = '<i class="fa-solid fa-battery-' + bc + '"></i><div class="i-val">' + ba + '<span>%</span></div>';

    this.ref_info(el, msg);
  }

  ref_dista_cont(some) {
    if (!some.ros) return;
    if (!some.ros.ds) return;

    let sid = 'dista-' + some.cid;
    let el = document.getElementById(sid);
    if (!el) return;

    const dis = this.fm_distance(some.ros.ds);

    let msg = '';
    if (dis.length > 0) {
      msg = '<i class="fa-solid fa-compass-drafting"></i><div class="i-val">' + dis + '</div>';
    }

    this.ref_info(el, msg);
  }

  ref_cnt() {
    this.user_cnt.innerHTML = '';
    let cc = this.list.children.length;
    let cc_inse = this.list.querySelectorAll('li.in-se').length;

    if (cc == 0) return;

    let msg = cc;
    if (parseInt(cc_inse) > 0) {
      msg = cc_inse + '/' + cc;
    }

    this.user_cnt.textContent = msg;
  }

  sync_litems() {
    const cids = Object.keys(this.oin.uslist);

    if (cids.length == 0) return;

    cids.forEach((cid, _) => {
      let litem = document.getElementById(cid);
      const some = this.oin.uslist[cid];

      if (some.in_se) {
        litem.classList.add('in-se');
      } else {
        litem.classList.remove('in-se');
      }

      if (some.some_se) {
        litem.classList.add('some-se');
      } else {
        litem.classList.remove('some-se');
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

      if (some.roomid) {
        litem.classList.add('w-room');
      } else {
        litem.classList.remove('w-room');
      }

      if (some.in_mon) {
        litem.classList.add('in-mon');
      } else {
        litem.classList.remove('in-mon');
      }

      this.ref_dista_cont(some);
    });

    this.ref_cnt();
  }

  is_located(some) {
    if (some.in_se) return;

    some.some_se = true;
    this.sync_litems();

    if (some.loc_tm) {
      clearTimeout(some.loc_tm);
      some.loc_tm = null;
    }

    some.loc_tm = setTimeout(() => {
      some.some_se = false;
      this.sync_litems();
    }, this.time_out*2);
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
    }, this.time_out);
  }

  get_cust_icon(some) {
    if (some.nik.toLowerCase().includes('как')) {
      return '/static/images/map/ka.png';
    }

    return '';
  }

  ref_ma(cid) {
    let some = this.oin.uslist[cid];

    if (!some) return false;

    this.ref_coo_cont(some);
    this.ref_bat_cont(some);
    this.update_timer(some);

    if (!some.pos) return false;

    let sp = [some.pos.lat, some.pos.lng];
    let route_opened = this.oin.get_route_cid().length > 0;

    if (some.ma) {
      some.ma.setLatLng(sp);
    } else {
      let pop = some.nik ? some.nik : some.cid;
      let cicon = this.get_cust_icon(some);
      let ico = this.oin.mai;
      if (cicon.length > 0) {
        ico = L.icon({
          'iconUrl': cicon,
          'iconAnchor': [15, 30],
          'tooltipAnchor': [0, -30]
        });
      }
      some.ma = L.marker(sp, {
        'icon': ico
      }).addTo(this.oin.map);

      some.ma.bindTooltip(pop, {
        'className': 'some-tooltip',
        'direction': 'top'
      });
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

    if (!route_opened) {
      this.oin.move_to_ma(some);
    }

    this.is_located(some);
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

    some.in_se = true;
    some.some_se = false;
    this.sync_litems();

    this.oin.req_loca(some);

    return false;
  }

  us_route_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget;
    let li_cont = this.fun.parent(el, '.map-us-li');
    let cid = li_cont.getAttribute('id');
    const rcid = this.oin.get_route_cid();
    let in_route = cid == rcid;
    let some = this.oin.uslist[cid];

    if (!some.pos) {
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

  us_monit_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget;
    let li_cont = this.fun.parent(el, '.map-us-li');
    let cid = li_cont.getAttribute('id');
    let in_mon = li_cont.classList.contains('in-mon');
    let some = this.oin.uslist[cid];

    some.in_mon = !in_mon;
    this.sync_litems();

    this.oin.move_to_ma(some);

    return false;
  }

  us_chat_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget;
    let li_cont = this.fun.parent(el, '.map-us-li');
    let cid = li_cont.getAttribute('id');
    let some = this.oin.uslist[cid];

    if (!some) return;

    this.oin.req_chat(some);

    return false;
  }

  get_li_tag() {
    let ret = '<li id="#CID#" class="map-us-li">' +
      '<div class="map-us-cont">' +
      '<div id="nik-#CID#" class="nik-cont"></div>' +
      '<div class="info-cont">' +
      '<div class="info-cont-items">' +
      '<div class="coo-cont" id="bat-#CID#" title="Battery"></div>' +
      '<div class="coo-cont" id="coo-#CID#" title="Location"></div>' +
      '<div class="coo-cont" id="dista-#CID#" title="Distance"></div>' +
      '</div>' +
      '</div>' +
      '</div>' +
      '<div class="us-btn-cont">' +
      '<span class="us-monit" title="Moving map">' +
      '<i class="fa-solid fa-crosshairs"></i>' +
      '</span>' +
      '<span class="us-route" title="Make a route">' +
      '<i class="fa-solid fa-car-side"></i>' +
      '</span>' +
      '<span class="bet-us-route" title="Make a route between people">' +
      '<i class="fa-solid fa-person-walking-arrow-loop-left"></i>' +
      '</span>' +
      '<span class="us-chat" title="Join to chat">' +
      '<i class="fa-solid fa-microphone"></i>' +
      '</span>' +
      '</div>' +
      '</li>';

    return ret;
  }

  get_new_li(some) {
    let cid = some.cid;
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
        nel = this.get_new_li(some);
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
