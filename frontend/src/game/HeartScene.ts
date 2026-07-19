import Phaser from 'phaser';

const HEART_TINTS = [
  0xff4466, 0xff6b8a, 0xff2244, 0xff8899, 0xff5577, 0xee3355,
];

export default class HeartScene extends Phaser.Scene {
  constructor() {
    super({ key: 'HeartScene' });
  }

  create() {
    this.generateHeartTexture();
  }

  private generateHeartTexture() {
    const size = 32;
    const gfx = this.add.graphics();

    gfx.fillStyle(0xff4466, 1);

    // Heart shape using two circles and a triangle
    const cx = size / 2;
    const cy = size / 2;
    const r = size * 0.22;

    gfx.fillCircle(cx - r * 0.85, cy - r * 0.3, r);
    gfx.fillCircle(cx + r * 0.85, cy - r * 0.3, r);

    // Bottom triangle
    gfx.beginPath();
    gfx.moveTo(cx - r * 1.85, cy - r * 0.1);
    gfx.lineTo(cx, cy + r * 1.6);
    gfx.lineTo(cx + r * 1.85, cy - r * 0.1);
    gfx.closePath();
    gfx.fillPath();

    gfx.generateTexture('heart', size, size);
    gfx.destroy();
  }

  spawnHeart(screenX: number, screenY: number) {
    const tint = HEART_TINTS[Math.floor(Math.random() * HEART_TINTS.length)];
    const scale = 0.5 + Math.random() * 0.5;
    const drift = (Math.random() - 0.5) * 60;
    const rise = 120 + Math.random() * 80;
    const duration = 800 + Math.random() * 600;

    const heart = this.add
      .image(screenX, screenY, 'heart')
      .setScale(scale)
      .setTint(tint)
      .setAlpha(1);

    this.tweens.add({
      targets: heart,
      x: screenX + drift,
      y: screenY - rise,
      alpha: 0,
      scale: scale * 0.4,
      duration,
      ease: 'Cubic.easeOut',
      onComplete: () => heart.destroy(),
    });
  }
}
