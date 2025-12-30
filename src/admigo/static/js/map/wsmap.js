;"use strict";
class Wsmap {
  constructor(oin_in) {
    this.oin = oin_in;

    this.ws = {
      'wsurl': '',
      'cid': '',
      'startpoint': '',
      'routeurl': '',
      'routekey': '',
      'handler': null,
      'TPS': {
        'RLOCA': 'rloca',
        'ALOCA': 'aloca',
        'CLIST': 'clist',
        'SENDERHI': 'sender_hi',
        'SENDERST': 'sender_stop',
        'GOCHAT': 'go_chat'
      }
    };

    this.initConf()
  }

  initConf() {
    const cs = document.getElementsByName('gorilla.csrf.Token')[0].value;
    const url = this.oin.elmap.getAttribute('href');

    axios.post(url, null, {
      headers: { 'X-CSRF-Token': cs }
    })
      .then((re) => {
        this.ws.cid = re.data.cid;
        this.ws.wsurl = re.data.link;
        this.ws.startpoint = re.data.startpoint;
        this.ws.routeurl = re.data.routeurl;
        this.ws.routekey = re.data.routekey;
        this.oin.init_map();
      })
      .catch(err => {
        console.log(err);
        this.oin.showLog(err, true);
      });
  }

  wsClear() {
    this.ws.handler = null;
  }

  wsError(ev) {
    this.oin.showLog('WS error: ' + ev.target.url, true);
  }

  wsOpen() {
    this.oin.showLog('Connected to server', false);
    this.clients_list();
  }

  wsClose(e) {
    this.oin.showLog('wsClose code: ' + e.code, true);
    this.wsClear()
  }

  wsMessage(e) {
    let jsi = e.data.split("\n");

    for (let i in jsi) {
      let msg = JSON.parse(jsi[i]);

      this.oin.showLog('wsMessage: ' + msg.tp, false);

      switch (msg.tp) {
        case this.ws.TPS.ALOCA:
          this.oin.ans_loca(msg);
          break;
        case this.ws.TPS.CLIST:
          this.oin.ref_uslist(msg.content);
          break;
        case this.ws.TPS.SENDERHI:
          this.oin.sender_hi(msg);
          break;
        case this.ws.TPS.SENDERST:
          this.oin.sender_stop(msg);
          break
      }
    }
  }

  startWs() {
    this.ws.handler = new WebSocket(this.ws.wsurl);
    this.ws.handler.onerror = this.wsError.bind(this);
    this.ws.handler.onopen = this.wsOpen.bind(this);
    this.ws.handler.onclose = this.wsClose.bind(this);
    this.ws.handler.onmessage = this.wsMessage.bind(this);
  }

  clients_list() {
    if (!this.ws.handler) return;

    let jo = {
      'tp': this.ws.TPS.CLIST,
      'content': ''
    };

    this.ws.handler.send(JSON.stringify(jo));
  }

  req_loca_ws(some) {
    if (!this.ws.handler) return;

    let jo = {
      'tp': this.ws.TPS.RLOCA,
      'cid': some.cid
    };

    this.ws.handler.send(JSON.stringify(jo));
  }

  req_chat_ws(some) {
    if (!this.ws.handler) return;

    if (!some.roomid) {
      some.roomid = this.oin.make_roomid();
    }

    let jo = {
      'tp': this.ws.TPS.GOCHAT,
      'cid': some.cid,
      'roomid': some.roomid,
    };

    this.ws.handler.send(JSON.stringify(jo));

    this.oin.sync_litems();
    let buf = 'meet/' + some.roomid;
    this.oin.cp_into_buf(buf);

    if (window.AndroidChatInterface) {
      window.AndroidChatInterface.sendRoomid(some.roomid);
    }
  }
}
