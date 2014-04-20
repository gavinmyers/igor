import QtQuick 2.0
import GoExtensions 1.0
import QtGraphicalEffects 1.0


Rectangle {
    Game {
      id: game
      x: 0
      y: 0
      Timer {
         interval: 500; running: true; repeat: true
         onTriggered: game.update() 
      }
      Timer {
         interval: 500; running: true; repeat: false 
         onTriggered: game.init() 
      }

    }

    width: 1260 
    height: 960 
    color: "gray"
    Player {
      id: player
      x: 640
      y: 480
      z:999
      Timer {
         interval: 50; running: true; repeat: true
         onTriggered: player.update() 
      }
      Item {
        id: photo 
        width: 16
        height: 18
        clip: true
        Image {
            x: 0 
            y: -17 
            source: "DawnLike_1/Characters/Player0.png" 
        }
      }
    }
    MouseArea {
        anchors.fill: parent
        onClicked: player.handleClick(mouse.x, mouse.y)
    }

    property var floor_0_0: Component {
      Rectangle {
        z:0
        width:16
        height:18
        clip: true
        Image {
          x: -32 
          y: -128
          source: "DawnLike_1/Objects/Floor.png" 
        }
      }
    }
    property var floor_0_1: Component {
      Rectangle {
        z:0
        width:16
        height:18
        clip: true
        Image {
          x: -48 
          y: -128
          source: "DawnLike_1/Objects/Floor.png" 
        }
      }
    }


}


