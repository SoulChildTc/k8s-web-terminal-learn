<script setup lang="ts">
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { SearchAddon } from 'xterm-addon-search';
// import { AttachAddon } from 'xterm-addon-attach'
// import SockJS from  'sockjs-client';  
import SockJS from  'sockjs-client/dist/sockjs.min.js';  
import axios from "axios";
import { ref,reactive, onMounted } from 'vue';


let sock: any
let term: any
let searchAddon:any

async function InitSock(){
  let sessionid = await getSessionID()
  if (! sessionid) {
    console.log("连接建立失败,sessuionID获取失败")
    return
  }


  sock = new SockJS('http://127.0.0.1:8090/api/sockjs?');

  sock.onopen = function() {
      console.log('connection open');
      sock.send('{"Op":"bind","SessionID":"' + sessionid + '"}')
      initTerm()
      onTerminalResize()
  };

  sock.onmessage = function(e:any) {
      let msg = JSON.parse(e.data)
      term.write(msg.Data)
  };

  sock.onclose = function() {
      console.log('connection closed');
  };

}

async function getSessionID(){
  try {
    const response = await axios.get("http://127.0.0.1:8090/api/v1/pod/default/busybox-deployment-dcb89bc87-hrspk/shell/busybox?shell=sh")
    return response.data.data.id
  } catch (err) {
    console.error(err);
  }
}

function onTerminalSendString(str: string): void {
    if (sock) {
      sock.send(
        JSON.stringify({
          Op: 'stdin',
          Data: str,
          Cols: term.cols,
          Rows: term.rows,
        })
      );
    }else {
      console.log("连接已断开,请重新建立连接")
    }
}

function onTerminalResize(): void {
    if (sock) {
      sock.send(
        JSON.stringify({
          Op: 'resize',
          Cols: term.cols,
          Rows: term.rows,
        })
      );
    }else {
      console.log("连接已断开,请重新建立连接")
    }
}

function initTerm() {
  const xtermjsTheme = {
  cursorAccent: '#2D2E2C',
  foreground: '#F8F8F8',
  background: '#2D2E2C',
  selectionBackground: '#5DA5D533',
  black: '#1E1E1D',
  brightBlack: '#262625',
  red: '#CE5C5C',
  brightRed: '#FF7272',
  green: '#5BCC5B',
  brightGreen: '#72FF72',
  yellow: '#CCCC5B',
  brightYellow: '#FFFF72',
  blue: '#5D5DD3',
  brightBlue: '#7279FF',
  magenta: '#BC5ED1',
  brightMagenta: '#E572FF',
  cyan: '#5DA5D5',
  brightCyan: '#72F0FF',
  white: '#F8F8F8',
  brightWhite: '#FFFFFF'
};
  term = new Terminal({
      fontSize: 16,
      fontFamily: 'Monaco, Consolas, "Lucida Console", monospace',
      cursorBlink: true,
      scrollback: 1000,  // 滚动缓冲区大小, 
      // cursorStyle: 'bar',
      // cursorWidth: 7,
      theme: xtermjsTheme,
    });

  const fitAddon = new FitAddon();
  searchAddon = new SearchAddon();
  term.loadAddon(searchAddon);
  term.loadAddon(fitAddon);
  term.open(document.getElementById('xterm'));
  fitAddon.fit();
  term.onResize((size: any) => {
    onTerminalResize()
  })
  window.addEventListener('resize', () => fitAddon.fit());
  term.focus();
  
  term.write("Welcome to Web Console By \x1b[1;1;32mSoulChild!\r\n\x1b[0m")


  term.onData((data: any) => {
    onTerminalSendString(data)
  });
  // term.onKey((e:any)=>{
  //   if (e.key === "\r") {
  //     onTerminalSendString("\r")
  //   }
  // })
}

onMounted(async ()=>{
  InitSock()
})

function search(){
  searchAddon.findNext('Soul')
  // searchAddon.findPrevious('Soul')
}


</script>

<template>
  <div id="xterm" class="xterm" />
  <!-- <button @click="search">查找</button> -->
</template>

<style scoped>
/* header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
} */
</style>
