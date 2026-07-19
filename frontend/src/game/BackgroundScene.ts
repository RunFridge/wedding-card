import Phaser from 'phaser';
import eyeOpenIdle from '../assets/characters/eye_open_idle.png';
import eyeCloseIdle from '../assets/characters/eye_close_idle.png';
import eyeCloseSmiling from '../assets/characters/eye_close_smiling.png';
import chickenImg from '../assets/farm/chicken.png';
import chickImg from '../assets/farm/chick.png';

import {
  type Palette,
  type SceneContext,
  type SkyGradient,
  DAY,
  NIGHT,
  SUNSET_SKY,
  DAWN_SKY,
} from './draw/palettes';
import {
  drawSky,
  drawStars,
  drawMoon,
  drawClouds,
  drawSun,
  drawBirds,
} from './draw/sky';
import {
  drawDistantTrees,
  drawHills,
  drawGround,
  drawPath,
  drawFence,
} from './draw/terrain';
import {
  drawTrees,
  drawBushes,
  drawMushrooms,
  drawTallGrass,
  drawFlowers,
  drawFenceVines,
} from './draw/vegetation';
import {
  drawFireflies,
  drawButterflies,
  drawAnimals,
  drawCharacters,
} from './draw/creatures';

export default class BackgroundScene extends Phaser.Scene implements SceneContext {
  p: Palette;
  _dark: boolean;
  _transitioning: boolean;
  _pendingResize: Phaser.Structs.Size | null;
  _reducedMotion: boolean;
  _lastWidth!: number;
  _animals!: Array<{
    sprite: Phaser.GameObjects.Sprite;
    wanderTween: Phaser.Tweens.Tween | null;
  }>;
  _treeGfx!: Phaser.GameObjects.Graphics[];
  _onDocumentClick: (() => void) | null = null;
  _onThemeChange: ((e: Event) => void) | null = null;
  flowers!: Phaser.GameObjects.Graphics[];
  smileTimer: Phaser.Time.TimerEvent | null = null;

  constructor() {
    super({ key: 'BackgroundScene' });
    this._dark = document.documentElement.classList.contains('dark');
    this.p = this._dark ? NIGHT : DAY;
    this._transitioning = false;
    this._pendingResize = null;
    this._reducedMotion = window.matchMedia(
      '(prefers-reduced-motion: reduce)',
    ).matches;
  }

  preload(): void {
    this.load.image('char_open', eyeOpenIdle);
    this.load.image('char_blink', eyeCloseIdle);
    this.load.image('char_smile', eyeCloseSmiling);
    this.load.spritesheet('chicken', chickenImg, {
      frameWidth: 16,
      frameHeight: 16,
    });
    this.load.spritesheet('chick', chickImg, {
      frameWidth: 16,
      frameHeight: 16,
    });
  }

  addTween(
    config: Phaser.Types.Tweens.TweenBuilderConfig,
  ): Phaser.Tweens.Tween | null {
    if (this._reducedMotion) return null;
    return this.tweens.add(config);
  }

  addTimedEvent(
    config: Phaser.Types.Time.TimerEventConfig,
  ): Phaser.Time.TimerEvent | null {
    if (this._reducedMotion) return null;
    return this.time.addEvent(config);
  }

  addDelayedCall(
    delay: number,
    callback: () => void,
  ): Phaser.Time.TimerEvent | null {
    if (this._reducedMotion) return null;
    return this.time.delayedCall(delay, callback);
  }

  create(): void {
    this._lastWidth = this.cameras.main.width;
    this._animals = [];

    this.anims.create({
      key: 'chicken_idle',
      frames: this.anims.generateFrameNumbers('chicken', {
        start: 0,
        end: 3,
      }),
      frameRate: 4,
      repeat: -1,
    });
    this.anims.create({
      key: 'chicken_walk',
      frames: this.anims.generateFrameNumbers('chicken', {
        start: 4,
        end: 7,
      }),
      frameRate: 6,
      repeat: -1,
    });
    this.anims.create({
      key: 'chick_idle',
      frames: this.anims.generateFrameNumbers('chick', { start: 0, end: 3 }),
      frameRate: 4,
      repeat: -1,
    });
    this.anims.create({
      key: 'chick_walk',
      frames: this.anims.generateFrameNumbers('chick', { start: 4, end: 7 }),
      frameRate: 6,
      repeat: -1,
    });

    this.drawAll();
    this.scale.on('resize', this.handleResize, this);

    this._onThemeChange = (e: Event) => {
      if (this._transitioning) return;

      const detail = (e as CustomEvent<{ dark: boolean; duration?: number }>)
        .detail;
      const newDark = detail.dark;
      const duration = detail.duration || 1500;
      const half = duration / 2;

      this._transitioning = true;

      const w = this.cameras.main.width;
      const h = this.cameras.main.height;
      const skyH = h * 0.55;

      const oldScene = this.children.getAll().slice();
      const oldTrees = [...(this._treeGfx || [])];
      oldTrees.forEach((g) => g.setDepth(2500));

      const fromSky: SkyGradient = newDark
        ? { top: [...DAY.skyTop], bot: [...DAY.skyBot] }
        : { top: [...NIGHT.skyTop], bot: [...NIGHT.skyBot] };
      const midSky: SkyGradient = newDark ? SUNSET_SKY : DAWN_SKY;
      const toSky: SkyGradient = newDark
        ? { top: [...NIGHT.skyTop], bot: [...NIGHT.skyBot] }
        : { top: [...DAY.skyTop], bot: [...DAY.skyBot] };

      const skyLayer = this.add.graphics().setDepth(2000);
      const bodyLayer = this.add.graphics().setDepth(2001);
      const tintLayer = this.add.graphics().setDepth(1000);
      const tintHex = newDark ? 0xff8c32 : 0x6488c8;

      const sunX = w * 0.82;
      const sunRestY = h * 0.06;
      const moonX = w * 0.8;
      const moonRestY = h * 0.08;
      const horizonY = h * 0.52;

      const transStars = Array.from({ length: 40 }, (_, i) => ({
        x: Math.round((Math.sin(i * 127.1 + 0.5) * 0.5 + 0.5) * w),
        y: Math.round(
          (Math.cos(i * 311.7 + 0.5) * 0.5 + 0.5) * skyH * 0.85,
        ),
        brightness: 0.3 + (Math.sin(i * 47.3) * 0.5 + 0.5) * 0.7,
        size: i % 3 === 0 ? 2 : 1,
      }));

      const lerpArr = (
        a: [number, number, number],
        b: [number, number, number],
        t: number,
      ): [number, number, number] =>
        a.map((v, i) => Math.round(v + (b[i] - v) * t)) as [
          number,
          number,
          number,
        ];

      const renderSky = (
        top: [number, number, number],
        bot: [number, number, number],
        starAlpha: number,
      ): void => {
        skyLayer.clear();
        const steps = 32;
        for (let i = 0; i < steps; i++) {
          const f = i / steps;
          const r = Math.round(top[0] + f * (bot[0] - top[0]));
          const g = Math.round(top[1] + f * (bot[1] - top[1]));
          const b = Math.round(top[2] + f * (bot[2] - top[2]));
          skyLayer.fillStyle((r << 16) | (g << 8) | b, 1);
          skyLayer.fillRect(0, f * skyH, w, skyH / steps + 1);
        }
        if (starAlpha > 0.01) {
          transStars.forEach((s) => {
            skyLayer.fillStyle(0xffffff, s.brightness * starAlpha);
            skyLayer.fillRect(s.x, s.y, s.size, s.size);
          });
        }
      };

      const renderSun = (x: number, y: number, a: number): void => {
        if (a < 0.01) return;
        const near = Math.max(0, (y - h * 0.1) / (h * 0.45));
        bodyLayer.fillStyle(0xfff176, 0.25 * a * (1 + near));
        bodyLayer.fillCircle(x, y, 32 + near * 16);
        bodyLayer.fillStyle(0xffee58, 0.45 * a);
        bodyLayer.fillCircle(x, y, 22);
        bodyLayer.fillStyle(0xffeb3b, 0.9 * a);
        bodyLayer.fillCircle(x, y, 14);
      };

      const renderMoon = (x: number, y: number, a: number): void => {
        if (a < 0.01) return;
        bodyLayer.fillStyle(0xeeeedd, 0.08 * a);
        bodyLayer.fillCircle(x, y, 40);
        bodyLayer.fillStyle(0xeeeedd, 0.15 * a);
        bodyLayer.fillCircle(x, y, 28);
        bodyLayer.fillStyle(0xeeeedd, 0.9 * a);
        bodyLayer.fillCircle(x, y, 16);
        if (a > 0.3) {
          bodyLayer.fillStyle(0x1a2450, Math.min(1, a * 1.5));
          bodyLayer.fillCircle(x + 6, y - 3, 13);
        }
        bodyLayer.fillStyle(0xddddcc, 0.3 * a);
        bodyLayer.fillCircle(x - 5, y + 2, 3);
        bodyLayer.fillCircle(x - 2, y - 4, 2);
      };

      const renderTint = (a: number): void => {
        tintLayer.clear();
        if (a < 0.01) return;
        tintLayer.fillStyle(tintHex, 0.15 * a);
        tintLayer.fillRect(0, skyH, w, h - skyH);
        tintLayer.fillStyle(tintHex, 0.08 * a);
        tintLayer.fillRect(0, skyH - 25, w, 25);
      };

      const progress = { t: 0 };

      this.tweens.add({
        targets: progress,
        t: 1,
        duration: half,
        ease: 'Sine.easeInOut',
        onUpdate: () => {
          const t = progress.t;
          const starAlpha = newDark ? 0 : 1 - t;
          renderSky(
            lerpArr(fromSky.top, midSky.top, t),
            lerpArr(fromSky.bot, midSky.bot, t),
            starAlpha,
          );
          renderTint(t);
          bodyLayer.clear();
          if (newDark) {
            renderSun(sunX, sunRestY + (horizonY - sunRestY) * t, 1 - t);
          } else {
            renderMoon(
              moonX,
              moonRestY + (horizonY - moonRestY) * t,
              1 - t,
            );
          }
        },
        onComplete: () => {
          this._dark = newDark;
          this.p = newDark ? NIGHT : DAY;
          this.cameras.main.setBackgroundColor(newDark ? '#0f1f0a' : '#3a7d22');

          this.tweens.killAll();
          this.time.removeAllEvents();
          this.removeDocumentListener();
          this.drawAll();

          const overlays = new Set([skyLayer, bodyLayer, tintLayer]);
          const oldSet = new Set(oldScene);
          const newScene = this.children
            .getAll()
            .filter(
              (c: Phaser.GameObjects.GameObject) =>
                !overlays.has(c as Phaser.GameObjects.Graphics) &&
                !oldSet.has(c),
            );

          newScene.forEach((c) => (c as unknown as Phaser.GameObjects.Components.Alpha).setAlpha(0));
          const newTrees = [...(this._treeGfx || [])];
          newTrees.forEach((g) => g.setDepth(2500));

          const p2 = { t: 0 };
          this.tweens.add({
            targets: p2,
            t: 1,
            duration: half,
            ease: 'Sine.easeInOut',
            onUpdate: () => {
              const t = p2.t;
              const starAlpha = newDark ? t : 0;
              renderSky(
                lerpArr(midSky.top, toSky.top, t),
                lerpArr(midSky.bot, toSky.bot, t),
                starAlpha,
              );
              renderTint(1 - t);
              bodyLayer.clear();
              if (newDark) {
                renderMoon(
                  moonX,
                  horizonY + (moonRestY - horizonY) * t,
                  t,
                );
              } else {
                renderSun(sunX, horizonY + (sunRestY - horizonY) * t, t);
              }
              oldScene.forEach((c) => (c as unknown as Phaser.GameObjects.Components.Alpha).setAlpha(1 - t));
              newScene.forEach((c) => (c as unknown as Phaser.GameObjects.Components.Alpha).setAlpha(t));
            },
            onComplete: () => {
              skyLayer.destroy();
              bodyLayer.destroy();
              tintLayer.destroy();
              oldScene.forEach((c) => c.destroy());
              newScene.forEach((c) => (c as unknown as Phaser.GameObjects.Components.Alpha).setAlpha(1));
              newTrees.forEach((g) => g.setDepth(0));
              this._transitioning = false;

              if (this._pendingResize) {
                const size = this._pendingResize;
                this._pendingResize = null;
                this.handleResize(size);
              }
            },
          });
        },
      });
    };
    window.addEventListener('theme-change', this._onThemeChange);
  }

  handleResize(gameSize: Phaser.Structs.Size): void {
    if (this._transitioning) {
      this._pendingResize = gameSize;
      return;
    }

    if (gameSize.width === this._lastWidth) {
      this.cameras.main.setSize(gameSize.width, gameSize.height);
      return;
    }
    this._lastWidth = gameSize.width;
    this.cameras.main.setSize(gameSize.width, gameSize.height);
    this.children.removeAll(true);
    this.tweens.killAll();
    this.time.removeAllEvents();
    this.removeDocumentListener();
    this.drawAll();
  }

  removeDocumentListener(): void {
    if (this._onDocumentClick) {
      document.removeEventListener('pointerdown', this._onDocumentClick);
      this._onDocumentClick = null;
    }
  }

  shutdown(): void {
    this.removeDocumentListener();
    if (this._onThemeChange) {
      window.removeEventListener('theme-change', this._onThemeChange);
      this._onThemeChange = null;
    }
  }

  destroy(): void {
    this.removeDocumentListener();
    if (this._onThemeChange) {
      window.removeEventListener('theme-change', this._onThemeChange);
      this._onThemeChange = null;
    }
  }

  drawAll(): void {
    const w = this.cameras.main.width;
    const h = this.cameras.main.height;

    drawSky(this, w, h);
    if (this._dark) {
      drawStars(this, w, h);
      drawMoon(this, w, h);
    } else {
      drawClouds(this, w, h);
      drawSun(this, w, h);
      drawBirds(this, w, h);
    }
    drawDistantTrees(this, w, h);
    drawHills(this, w, h);
    drawGround(this, w, h);
    drawPath(this, w, h);
    drawFence(this, w, h);
    drawTrees(this, w, h);
    drawBushes(this, w, h);
    drawMushrooms(this, w, h);
    drawTallGrass(this, w, h);
    drawFlowers(this, w, h);
    drawFenceVines(this, w, h);
    if (this._dark) {
      drawFireflies(this, w, h);
    } else {
      drawButterflies(this, w, h);
    }
    drawAnimals(this, w, h);
    drawCharacters(this, w, h);
  }
}
