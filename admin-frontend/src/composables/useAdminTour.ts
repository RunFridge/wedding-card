import { driver, type DriveStep } from 'driver.js';
import 'driver.js/dist/driver.css';
import { i18n } from '@/i18n';

const STORAGE_KEY = 'admin-tour-seen';

function buildSteps(): DriveStep[] {
  const { t } = i18n.global;
  const steps: DriveStep[] = [
    {
      popover: {
        title: t('tour.welcomeTitle'),
        description: t('tour.welcomeDescription'),
      },
    },
    {
      element: '[data-tour="nav-guestbook"]',
      popover: {
        title: t('tour.guestbookTitle'),
        description: t('tour.guestbookDescription'),
      },
    },
    {
      element: '[data-tour="nav-rankings"]',
      popover: {
        title: t('tour.rankingsTitle'),
        description: t('tour.rankingsDescription'),
      },
    },
    {
      element: '[data-tour="nav-hall-of-fame"]',
      popover: {
        title: t('tour.hallOfFameTitle'),
        description: t('tour.hallOfFameDescription'),
      },
    },
    {
      element: '[data-tour="nav-photos"]',
      popover: {
        title: t('tour.photosTitle'),
        description: t('tour.photosDescription'),
      },
    },
    {
      element: '[data-tour="nav-asset-photos"]',
      popover: {
        title: t('tour.assetPhotosTitle'),
        description: t('tour.assetPhotosDescription'),
      },
    },
    {
      element: '[data-tour="nav-settings"]',
      popover: {
        title: t('tour.settingsTitle'),
        description: t('tour.settingsDescription'),
      },
    },
    {
      element: '[data-tour="nav-system"]',
      popover: {
        title: t('tour.systemTitle'),
        description: t('tour.systemDescription'),
      },
    },
    {
      element: '[data-tour="visitors"]',
      popover: {
        title: t('tour.visitorsTitle'),
        description: t('tour.visitorsDescription'),
      },
    },
    {
      element: '[data-tour="help"]',
      popover: {
        title: t('tour.helpTitle'),
        description: t('tour.helpDescription'),
      },
    },
  ];
  return steps.filter(
    (step) => !step.element || document.querySelector(step.element as string),
  );
}

export function startTour() {
  const { t } = i18n.global;
  driver({
    showProgress: true,
    nextBtnText: t('tour.next'),
    prevBtnText: t('tour.prev'),
    doneBtnText: t('tour.done'),
    progressText: '{{current}} / {{total}}',
    steps: buildSteps(),
  }).drive();
}

export function startTourOnFirstVisit() {
  if (localStorage.getItem(STORAGE_KEY)) return;
  localStorage.setItem(STORAGE_KEY, 'true');
  startTour();
}
