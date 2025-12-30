;"use strict";
class Agrida {
  constructor(props_in) {
    this.props = props_in;

    this.gr_cont = this.props.gr_cont;
    this.gr_mn = null;
    this.rand = this.gr_cont.getAttribute('data-rand');

    this.can_normal_row = true;
    this.hi_style_id = 'trst-allert-' + this.rand;
    this.set_gr_mn();
  }

  set_gr_mn() {
    if (this.gr_mn) return;
    this.gr_mn = this.gr_cont.querySelector('.rwd-menu');
    if (!this.gr_mn) return;
    document.body.appendChild(this.gr_mn);
  }

  docon() {
    document.querySelectorAll('#' + this.gr_cont.id + ' .rwd-table tbody tr.cla').forEach(tr => {
      if (!this.props.fun.once(tr, 'row_click')) {
        tr.addEventListener('click', this.row_click.bind(this));
      }
    });
  }

  highlight_row(tbl_id, rid_in) {
    let hi = document.getElementById(this.hi_style_id);

    if (hi) return;

    let sel_el = '#' + tbl_id + ' tr[data-rid="' + rid_in + '"]';
    let st = document.createElement('style');
    st.id = this.hi_style_id;
    st.innerHTML = sel_el + '{background-color:#e5e5e5}';
    document.getElementsByTagName('head')[0].appendChild(st);
  }

  normal_row() {
    let hi = document.getElementById(this.hi_style_id);

    if (!hi) return;
    if (!this.can_normal_row) return;

    hi.remove();
  }

  fill_menu(row) {
    let lis = '';
    let ul = this.gr_mn.querySelector('ul');

    row.querySelectorAll('.act-item').forEach(el => {
      let str = `
        <li>
          <div class="act clearfix" data-u="#URL#" data-m="#MET#">
            #ICO#
            <div class="act-cap">
              <span>#SPAN#</span>
              #WHO#
            </div>
          </div>
        </li>
      `;

      let ico = el.querySelector('i');

      let who = el.getAttribute('data-who');
      if (who) {
        who = '<span>' + who + '</span>';
      } else who = '';

      let mi = str
        .replace(/#URL#/g, el.getAttribute('data-u'))
        .replace(/#MET#/g, el.getAttribute('data-m'))
        .replace(/#ICO#/g, ico.outerHTML)
        .replace(/#SPAN#/g, el.getAttribute('title'))
        .replace(/#WHO#/g, who);

      lis += mi;
    });

    ul.innerHTML = lis;

    ul.querySelectorAll('.act').forEach(el => {
      if (this.props.fun.once(el, 'act_click')) return;
      el.addEventListener('click', this.actClick.bind(this));
    });
  }

  remove_clicked() {
    this.gr_cont.querySelectorAll('.clicked').forEach((cl) => {
      cl.classList.remove('clicked');
    });
  }

  hide_menu() {
    this.gr_mn.classList.add('hide');
    this.remove_clicked();
    this.normal_row();
  }

  show_menu(row) {
    this.remove_clicked();

    let rid_in = row.getAttribute('data-rid');
    let tid = this.gr_cont.id;

    this.normal_row();
    this.highlight_row(tid, rid_in);

    if (this.gr_mn.data && this.gr_mn.data.tm) {
      clearTimeout(this.gr_mn.data.tm);
    }

    if (!this.gr_mn.data) {
      this.gr_mn.data = {
        'tm': null
      }
    }

    if (!this.props.fun.once(this.gr_mn, 'mouse_enter')) {
      this.gr_mn.addEventListener('mouseenter', () => {
        clearTimeout(this.gr_mn.data.tm);
      });
    }

    if (!this.props.fun.once(this.gr_mn, 'mouse_leave')) {
      this.gr_mn.addEventListener('mouseleave', () => {
        this.gr_mn.data.tm = setTimeout(() => {
          this.hide_menu();
        }, 100)
      });
    }

    this.gr_mn.classList.remove('hide');
  }

  pos_menu(ev) {
    const pgx = ev.pageX;
    const mow = this.gr_mn.offsetWidth;
    const bow = document.body.offsetWidth;
    const pgy = ev.pageY;
    const moh = this.gr_mn.offsetHeight;
    const boh = document.body.offsetHeight;

    let new_left = pgx;
    if (pgx > mow ) new_left = pgx - mow;
    if ((bow - pgx) > mow) new_left = pgx;

    let new_top = pgy;
    if (moh > (boh - pgy)) new_top = pgy - moh;

    this.gr_mn.style.top = new_top + 'px';
    this.gr_mn.style.left = new_left + 'px';
  }

  row_click(ev) {
    this.set_gr_mn();

    let row = ev.currentTarget;

    if (row.classList.contains('clicked')) {
      this.hide_menu();
      return;
    }

    this.fill_menu(row);
    this.show_menu(row);

    row.classList.add('clicked');

    this.pos_menu(ev);
  }

  act_del(btn) {
    this.can_normal_row = false;

    let url = btn.getAttribute('data-u');

    let dlg = new Dlg({
      'fun': this.props.fun
    });

    dlg.show({
      'header': window.lang.re('Confirmation'),
      'msg': window.lang.re('Confirm deleting the row'),
      'buttons': [
        {'cap': window.lang.re('Cancel'), 'ret': 0},
        {'cap': window.lang.re('Delete'), 'ret': 1}
      ]
    })
      .then((re) => {
        if (re <= 0) return null;

        return axios.delete(url, {
          headers: {'X-CSRF-Token': this.props.cs}
        });
      })
      .then((ret) => {
        if (ret == null) return;
        window.location.reload();
      })
      .catch((err) => {
        if (err.response.data.errors.api) {
          return dlg.show({
            'header': 'Error',
            'msg': err.response.data.errors.api,
            'buttons': [
              {'cap': 'Ok', 'ret': 1}
            ]
          });
        }
      })
      .finally(() => {
        this.can_normal_row = true;
        this.normal_row();
      });
  }

  act_get(btn) {
    let url = btn.getAttribute('data-u');
    window.location.href = url;
  }


  actClick(ev) {
    ev.preventDefault();
    ev.stopPropagation();

    let btn = ev.currentTarget;
    let met = btn.getAttribute('data-m');

    switch (met) {
      case 'delete':
        this.act_del(btn);
        break;
      case 'get':
        this.act_get(btn);
        break;
      default:
        console.log('unknown method');
    }

    return false;
  }
}
