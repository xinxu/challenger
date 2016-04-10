import React from 'react';
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/board.css'

const Board = CSSModules(observer(React.createClass({
  render() {
    let match = this.props.game.match
    return (
      <div>
      <table>
      <tbody>
      <tr>
      <th>游戏模式</th>
      <th>游戏人数</th>
      <th>{match.mode == 1 ? '金币总数' : '游戏时长'}</th>
      <th>暴走次数</th>
      </tr>
      <tr>
      <td>{match.mode == 1 ? '赏金模式' : '生存模式'}</td>
      <td>{match.member.length}</td>
      <td>{match.mode == 1 ? match.gold : match.elasped}</td>
      <td>{match.rampageCount}</td>
      </tr>
      </tbody>
      </table>
      <table styleName='player'>
      <tbody>
      <tr>
      <th>玩家</th>
      <th>获得金币</th>
      <th>获得能量</th>
      <th>T0-T1</th>
      <th>T1-T2</th>
      <th>T2-T3</th>
      <th>大于T3</th>
      <th>被捕次数</th>
      </tr>
      {
        match.member.map((member) => {
          return (
            <tr key={member.name}>
            <td>{member.name}</td>
            <td>{member.gold}</td>
            <td>{member.energy}</td>
            <td>{member.levelData[0]}</td>
            <td>{member.levelData[1]}</td>
            <td>{member.levelData[2]}</td>
            <td>{member.levelData[3]}</td>
            <td>{member.hitCount}</td>
            </tr>)
        })
      }
      </tbody>
      </table>
      </div>
    );
  }
})), styles)

export default Board