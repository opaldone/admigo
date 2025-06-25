;"use strict";
class SeList {
  constructor(props_in) {
    this.props = props_in;

    this.form = null;
    if (this.props.able.getAttribute('form-submit') != null) {
      this.form = document.getElementById(this.props.able.getAttribute('form-submit'));
    }

    this.list_to_add = null;
    this.items_added = null;
    this.item_cuq = null;
    this.init_list_to_add();

    this.inp = this.props.able.querySelector('input[type=text]');
    this.se_list = this.props.able.querySelector('.se-list');

    this.inp.addEventListener('input', this.inp_input.bind(this));
    this.inp.addEventListener('keydown', this.inp_keydown.bind(this));
    this.inp.addEventListener('focusin', this.inp_focusin.bind(this));
    this.inp.addEventListener('focusout', this.inp_focusout.bind(this));

    this.se_list.addEventListener('focusout', this.se_list_focusout.bind(this));
  }

  mi_item_remove_click(ev) {
    let btn = ev.currentTarget;
    let par_ml_item = this.props.fun.parent(btn, '.ml-item');
    let li_val = btn.getAttribute('data-vuq');

    delete this.items_added[li_val];
    par_ml_item.remove();
  }

  set_mi_item_remove() {
    if (!this.list_to_add) return;
    this.list_to_add.querySelectorAll('.sp-del').forEach(btn => {
      if (this.props.fun.once(btn, 'rem_click')) return;
      if (!btn.getAttribute('data-vuq')) return;
      btn.addEventListener('click', this.mi_item_remove_click.bind(this));
    });
  }

  init_list_to_add() {
    if (!this.props.able.getAttribute('list-to-add')) return;

    this.items_added = {};
    this.list_to_add = document.getElementById(this.props.able.getAttribute('list-to-add'));
    this.item_cuq = this.list_to_add.getAttribute('data-cuq');

    this.list_to_add.querySelectorAll('.ml-item').forEach(liem => {
      let item = liem.querySelector('input[type=hidden]').value
      this.items_added[item] = true;
    });

    this.set_mi_item_remove();
  }

  clear_search() {
    this.se_list.innerHTML = '';
  }

  get_lis() {
    const ul = this.se_list.querySelector('.search-list-ul');

    if (!ul) return null;

    const list = ul.querySelectorAll('li');
    const lise = ul.querySelector('li.sele');
    const ixsel = Array.prototype.indexOf.call(list, lise);

    return {
      'list': list,
      'lise': lise,
      'ixsel': ixsel,
      'ul': ul
    };
  }

  add_to_list(li) {
    const li_val = li.getAttribute('data-' + this.item_cuq);

    if (this.items_added[li_val]) return;

    const lid = this.list_to_add.id;

    const cp_item = document.getElementById('templ-' + lid);
    let stri = '<li class="ml-item">' +
      cp_item.innerHTML +
    '</li>';

    let tem = document.createElement('template');
    tem.innerHTML = stri;
    let tem_cont = tem.content;

    let inp_hid = tem_cont.querySelector('input[type=hidden]');
    inp_hid.name = lid + '[]';
    inp_hid.value = li_val;
    let sea_el = tem_cont.querySelector('.sea-tmp');
    sea_el.classList.remove('sea-tmp');

    let btn_rem = tem_cont.querySelector('.sp-del');
    btn_rem.setAttribute('data-vuq', li_val);

    this.items_added[li_val] = true;
    this.list_to_add.append(tem.content);

    this.set_mi_item_remove();
  }

  set_sea(li) {
    for (const [ke, va] of Object.entries(li.dataset)) {
      const cls = '.sea-' + ke;

      document.querySelectorAll(cls).forEach(ele_in => {
        if (ele_in.classList.contains('sea-tmp')) return;
        if (ele_in.classList.contains('sea-inhtml')) {
          if (ele_in.classList.contains('sea-ifemp') && ele_in.textContent.length > 0) return;
          ele_in.textContent = va;
          return;
        }

        if (ele_in.classList.contains('sea-ifemp') && ele_in.value.length > 0) return;
        ele_in.value = va;
        this.props.fun.trigger(ele_in, 'keyup');
      });
    }

    if (this.list_to_add) {
      this.inp.value = '';
    }
  }

  li_s_click(li) {
    if (this.list_to_add) {
      this.add_to_list(li);
    }

    this.set_sea(li);

    this.clear_search();

    if (this.form) {
      this.form.submit();
    }
  }

  docon() {
    let lis = this.get_lis();
    if (!lis) return;

    if (!lis.list) return;
    if (lis.list.length == 0) return;

    lis.list.forEach((lit) => {
      if (this.props.fun.once(lit, 'lit_click')) return;
      lit.addEventListener('click', (ev) => {
        ev.stopPropagation();
        ev.preventDefault();
        let ct = ev.currentTarget;
        this.li_s_click(ct);
        return false;
      });
    });
  }

  inp_input(ev) {
    this.clear_search();

    const ct = ev.currentTarget;
    const str = ct.value;

    if (str.length == 0) return;

    let url = this.se_list.getAttribute('data-seu');

    url = url + '?' + 'fi=' + str;

    axios.get(url)
      .then((re) => {
        if (!re.data || !re.data.cont) {
          this.clear_search();
          return;
        }

        this.se_list.innerHTML = re.data.cont;
        this.docon();
      })
      .catch((err) => {
        console.log(err);
      });
  }

  inp_arrow(lis, ev) {
    let sel = lis.lise;

    if (sel == null) {
      lis.list[0].classList.add('sele');
      return false;
    }

    let st = (ev.key == 'ArrowDown' ? 1 : -1);
    let ix_sel = lis.ixsel;
    let ix_nex = ix_sel + st;
    let len = lis.list.length;

    if (ix_nex >= len) return false;
    if (ix_nex < 0) return false;

    sel.classList.remove('sele');
    lis.list[ix_nex].classList.add('sele');

    lis.list[ix_nex].scrollIntoView({ 'behavior': 'smooth', 'block': 'end', 'inline': 'nearest' });

    return false;
  }

  inp_keydown(ev) {
    if (!['Enter', 'ArrowDown', 'ArrowUp'].includes(ev.key)) return true;

    ev.stopPropagation();
    ev.preventDefault();

    const lis = this.get_lis();

    if (!lis || lis.list.length == 0) {
      return false;
    }

    if (['ArrowDown', 'ArrowUp'].includes(ev.key)) {
      return this.inp_arrow(lis, ev);
    }

    if (lis.lise == null) {
      this.li_s_click(lis.list[0]);
      return false;
    }

    this.li_s_click(lis.lise);
    return false;
  }

  inp_focusin() {
    this.props.fun.trigger(this.inp, 'input');
  }

  inp_focusout(ev) {
    if (this.se_list == ev.relatedTarget) {
      return true;
    }

    this.clear_search();
  }

  se_list_focusout() {
    this.clear_search();
  }
}

class SeListHandler {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    document.querySelectorAll('.search-able').forEach(elem => {
      new SeList({
        'fun': this.fun,
        'able': elem
      })
    });
  }
}

new SeListHandler();
