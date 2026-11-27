class Dlg {
  constructor() {
    this.sh_cl = 'show-dlg-screen';
    this.b_cl = 'dlg-btn';

    this.scr = null;
    this.bd = null;
    this.header = null;
    this.msg = null;
    this.buttons = null;
    this.clobut = null;
    this.prom_good = null;
  }

  set_header(pars) {
    this.header.innerHTML = '';

    if (!pars.header) {
      this.header.classList.add('emp');
      return;
    }

    this.header.classList.remove('emp');
    this.header.textContent = pars.header;
  }

  set_msg(pars) {
    this.msg.innerHTML = '';

    if (!pars.msg) {
      this.msg.classList.add('emp');
      return;
    }

    this.msg.classList.remove('emp');
    this.msg.textContent = pars.msg;
  }

  btn_click(ev) {
    let th = ev.currentTarget;
    let ret = parseInt(th.getAttribute('data-ret'));
    this.close_dlg(ret);
  }

  set_buttons(pars) {
    this.buttons.innerHTML = '';

    if (!pars.buttons) {
      this.buttons.classList.add('emp');
      return;
    }

    this.buttons.classList.remove('emp');

    pars.buttons.forEach((btn) => {
      let b_el = document.createElement('button');
      b_el.classList.add(this.b_cl);
      if (btn.ret <= 0) b_el.classList.add('cancel');
      b_el.setAttribute('data-ret', btn.ret);
      b_el.textContent = btn.cap;
      b_el.addEventListener('click', this.btn_click.bind(this));
      this.buttons.appendChild(b_el);
    });
  }

  empty_prom() {
    return new Promise((good) => {
      good(0);
    });
  }

  clobut_click() {
    this.close_dlg();
  }

  window_esc(e) {
    if (e.key === 'Escape') {
      this.close_dlg();
    }
  }

  add_events() {
    window.addEventListener('keydown', this.window_esc);
    this.clobut.addEventListener('click', this.clobut_click);
  }

  remove_events() {
    window.removeEventListener('keydown', this.window_esc);
    this.clobut.removeEventListener('click', this.clobut_click);
  }

  set_elems() {
    this.scr = document.getElementById('dlg-screen');

    if (!this.scr) return false;

    this.bd = document.body;

    if (this.bd.classList.contains(this.sh_cl)) return false;

    this.header = document.getElementById('dlg-header');
    this.msg = document.getElementById('dlg-msg');
    this.buttons = document.getElementById('dlg-btns');
    this.clobut = document.getElementById('dlg-close');

    this.window_esc = this.window_esc.bind(this);
    this.clobut_click = this.clobut_click.bind(this);

    this.add_events();

    return true;
  }

  set_screen_class(pars) {
    this.scr.removeAttribute('class');

    if (pars.screen_class) {
      this.scr.classList.add(pars.screen_class);
    }
  }

  init_dlg(pars) {
    this.set_header(pars);
    this.set_msg(pars);
    this.set_buttons(pars);
    this.set_screen_class(pars);
  }

  close_dlg(ret_in) {
    let ret = 0;
    if (ret_in != undefined) ret = ret_in;

    this.remove_events();

    this.bd.classList.remove(this.sh_cl);
    this.prom_good(ret);
  }

  show(pars_in) {
    if (!this.set_elems()) {
      return this.empty_prom();
    }

    let pars = {};
    if (pars_in) pars = pars_in;

    this.init_dlg(pars);
    this.bd.classList.add(this.sh_cl);
    this.scr.focus();

    if (pars.buttons) {
      this.buttons.querySelector("." + this.b_cl).focus();
    }

    return new Promise((good) => {
      this.prom_good = good;
    });
  }
}
