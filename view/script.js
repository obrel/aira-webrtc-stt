'use strict';

const audio = document.querySelector('audio');
const startButton = document.getElementById('startButton');
const results = document.getElementById('results');
const constraints = window.constraints = {
  audio: true,
  video: false
};
const offerOptions = {
  offerToReceiveAudio: false,
  offerToReceiveVideo: false
};

let localStream;
let peerConnection;
let dataChannel;
let start = false;

startButton.addEventListener('click', async () => {
  if (!start) {
    await navigator.mediaDevices.getUserMedia(constraints).then(stream => {
      localStream = stream;
      localStream.oninactive = function() {
        console.log('Stream ended.');
      };

      audio.srcObject = stream
    }).catch(error => {
      console.error(error);
    });

    const audioTracks = localStream.getAudioTracks();

    if (audioTracks.length > 0) {
      console.log(`Using audio device: ${audioTracks[0].label}`);
    }

    peerConnection = new RTCPeerConnection({
      iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
    });
    peerConnection.addEventListener('icecandidate', event => onIceCandidate(event));
    peerConnection.addEventListener('iceconnectionstatechange', event => onIceStateChange(event));

    localStream.getAudioTracks().forEach(track => peerConnection.addTrack(track, localStream));

    dataChannel = peerConnection.createDataChannel('transcription', {
      ordered: true,
      protocol: 'tcp'
    });
    dataChannel.onmessage = event => {
      decodeDataChannelPayload(event.data).then(data => {
        const result = JSON.parse(data);
        const li = document.createElement('li');
        li.appendChild(document.createTextNode(`${result.text} (${result.confidence*100}%)`));
        results.appendChild(li);
      });
    };

    try {
      const offer = await peerConnection.createOffer(offerOptions);

      try {
        await peerConnection.setLocalDescription(offer);
      } catch (error) {
        onSetSessionDescriptionError(error);
      }
    } catch (error) {
      onCreateSessionDescriptionError(error);
    }

    startButton.textContent = 'Stop';
    start = true;
  } else {
    if (localStream) {
      localStream.getAudioTracks().forEach(track => track.stop());
      localStream = null;
    }

    peerConnection.close();
    peerConnection = null;

    startButton.textContent = 'Start';
    start = false;
  }
});

async function onIceCandidate(event) {
  if (!event.candidate) {
    const opts = {
      method: 'POST',
      body: JSON.stringify({
          offer: peerConnection.localDescription.sdp
      })
    };

    const resp = await fetch('http://localhost:4000/signaling', opts);
    const answer = await resp.json();

    try {
      await peerConnection.setRemoteDescription(new RTCSessionDescription({
        type: 'answer',
        sdp: answer.answer
      }));
    } catch (error) {
      onSetSessionDescriptionError(error);
    }
  }
}

function onIceStateChange(event) {
  if (peerConnection) {
    console.log('ICE state change event: ', event);
  }
}

function onCreateSessionDescriptionError(error) {
  console.log(`Failed to create session description: ${error.toString()}`);
}

function onSetSessionDescriptionError(error) {
  console.log(`Failed to set session description: ${error.toString()}`);
}

function decodeDataChannelPayload(data) {
  if (data instanceof ArrayBuffer) {
    const dec = new TextDecoder('utf-8');

    return Promise.resolve(dec.decode(data));
  } else if (data instanceof Blob) {
    const reader = new FileReader();
    const readPromise = new Promise((accept, reject) => {
      reader.onload = () => accept(reader.result);
      reader.onerror = reject;
    });

    reader.readAsText(data, 'utf-8');
    return readPromise;
  }
}