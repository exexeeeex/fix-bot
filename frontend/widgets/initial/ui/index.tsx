import { Button } from '@/components/ui/button';
import { AnimatePresence, motion } from 'motion/react';
import { InitialForm } from './form';
import { Icon } from '@/shared/components';
import { useInitialForm } from '../hooks';

const screenVariants = {
  initial: { opacity: 0, y: 12, filter: 'blur(4px)' },
  animate: { opacity: 1, y: 0, filter: 'blur(0px)' },
  exit: { opacity: 0, y: -12, filter: 'blur(4px)' },
};

const screenTransition = {
  duration: 0.3,
  ease: [0.22, 1, 0.36, 1] as const,
};

export const Initial = () => {
  const form = useInitialForm();
  const { isVisibleForm, selectRole } = form;

  return (
    <motion.section
      className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden bg-transparent px-6 pb-20 pt-12 text-center"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.4 }}
    >
      <AnimatePresence mode="wait" initial={false}>
        {!isVisibleForm ? (
          <motion.div
            key="role-selection"
            variants={screenVariants}
            initial="initial"
            animate="animate"
            exit="exit"
            transition={screenTransition}
            className="relative flex w-full max-w-md flex-col items-center"
          >
            <motion.div
              aria-hidden="true"
              className="pointer-events-none absolute -left-24 top-1/4 h-64 w-64 rounded-full bg-[#007aff]/10 blur-3xl"
              animate={{
                opacity: [0.35, 0.75, 0.35],
                scale: [0.9, 1.08, 0.9],
              }}
              transition={{
                duration: 6,
                repeat: Infinity,
                ease: 'easeInOut',
              }}
            />

            <motion.h1
              className="text-5xl font-extrabold tracking-tight text-primary/80 sm:text-6xl"
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{
                duration: 0.5,
                ease: [0.22, 1, 0.36, 1],
              }}
            >
              ФИКС
            </motion.h1>

            <motion.div
              className="mt-5 flex w-full flex-col gap-3"
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{
                duration: 0.5,
                delay: 0.1,
                ease: [0.22, 1, 0.36, 1],
              }}
            >
              <motion.div whileHover={{ y: -2 }} whileTap={{ scale: 0.98 }}>
                <Button
                  onClick={() => selectRole('student')}
                  className="h-14 w-full justify-between rounded-2xl px-5 text-base font-semibold shadow-lg shadow-primary/20 transition-shadow hover:shadow-xl hover:shadow-primary/25"
                >
                  <span className="flex items-center gap-3">
                    <Icon icon="graduation-cap" className="opacity-70" />Я студент
                  </span>

                  <Icon icon="arrow-right" className="opacity-70" />
                </Button>
              </motion.div>

              <motion.div whileHover={{ y: -2 }} whileTap={{ scale: 0.98 }}>
                <Button
                  onClick={() => selectRole('teacher')}
                  variant="outline"
                  className="h-14 w-full justify-between rounded-2xl bg-muted/20 px-5 text-base font-semibold shadow-sm backdrop-blur transition-shadow hover:bg-white hover:shadow-md"
                >
                  <span className="flex items-center gap-3">
                    <Icon icon="user" size={50} className="opacity-70" />Я преподаватель
                  </span>

                  <Icon icon="arrow-right" className="opacity-70" />
                </Button>
              </motion.div>
            </motion.div>
          </motion.div>
        ) : (
          <motion.div
            key="initial-form"
            variants={screenVariants}
            initial="initial"
            animate="animate"
            exit="exit"
            transition={screenTransition}
            className="w-full"
          >
            <InitialForm form={form} />
          </motion.div>
        )}
      </AnimatePresence>
    </motion.section>
  );
};
