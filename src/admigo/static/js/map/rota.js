;"use strict";
class Rota {
  constructor(oin_in) {
    this.oin = oin_in;
    this.press_dur = 100
    this.st = 3;

    document.querySelectorAll('.map-rot').forEach(el => {
      if (this.oin.fun.once(el, 'map_rot_tm')) return;
      el.data = {
        'press_tm': null
      };

      el.addEventListener('mousedown', this.start_click.bind(this));
      el.addEventListener('mouseup', this.end_click.bind(this));
      el.addEventListener('mouseleave', this.el_leave.bind(this));

      el.addEventListener('touchstart', this.start_click.bind(this), { passive: true });
      el.addEventListener('touchend', this.end_click.bind(this));
      el.addEventListener('touchcancel', this.el_leave.bind(this));
    });
  }

  rot_click(el) {
    let st = this.st;
    if (el.classList.contains('le')) st = -1 * st;
    let gb = this.oin.map.getBearing() + st;
    this.oin.map.setBearing(gb);
  }

  start_click(e) {
    let el = e.currentTarget;

    el.data.press_tm = setInterval(() => {
      this.rot_click(el);
    }, this.press_dur);
  }

  end_click(e) {
    let el = e.currentTarget;
    clearInterval(el.data.press_tm);
    this.rot_click(el);
  }

  el_leave(e) {
    let el = e.currentTarget;
    clearInterval(el.data.press_tm);
  }
}
