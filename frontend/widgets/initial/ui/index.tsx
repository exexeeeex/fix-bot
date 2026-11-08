import { Button } from '@/components/ui/button';
import { AnimatePresence, motion } from 'motion/react';
import { useState } from 'react';
import { useNavigate } from 'react-router';
import { ApplicationForm } from '../../../features/application/ui/form';
import { Icon } from '@/shared/components';
import { useDemoLogin } from '@/features/authentication/hooks';
import { useInitialForm } from '../hooks';
import { useApplicationForm } from '@/features/application/hooks';
import { AuthenticationStatusCard } from './application-status/authentication-status-card';
import { PendingApplicationCard } from './application-status/pending-application-card';
import { RejectedApplicationCard } from './application-status/rejected-application-card';

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
  const { userId, name, status, retryAuthentication } = useInitialForm();
  const applictionForm = useApplicationForm(userId, name);
  const { isVisibleForm, selectRole } = applictionForm;
  const navigate = useNavigate();
  const { login: demoLogin, isLoading: isDemoLoading } = useDemoLogin();
  const [demoError, setDemoError] = useState(false);

  const enterDemo = async (role: 'student' | 'teacher') => {
    setDemoError(false);
    try {
      await demoLogin(role);
      navigate(role === 'student' ? '/home' : '/lessons', { replace: true });
    } catch (error) {
      console.error('Failed to enter demo mode', error);
      setDemoError(true);
    }
  };

  return (
    <motion.section
      className="relative flex min-h-svh flex-col items-center justify-center overflow-x-hidden bg-transparent px-4 py-6 text-center sm:min-h-screen sm:overflow-visible sm:px-6 sm:py-12"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.4 }}
    >
      <AnimatePresence mode="wait" initial={false}>
        {status !== 'none' ? (
          <motion.div
            key={`application-status-${status}`}
            variants={screenVariants}
            initial="initial"
            animate="animate"
            exit="exit"
            transition={screenTransition}
            className="w-full max-w-md"
          >
            {status === 'pending' ? (
              <PendingApplicationCard />
            ) : status === 'rejected' ? (
              <RejectedApplicationCard />
            ) : (
              <AuthenticationStatusCard
                status={status}
                onRetry={() => void retryAuthentication()}
              />
            )}
          </motion.div>
        ) : !isVisibleForm ? (
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
            <ApplicationForm form={applictionForm} />
          </motion.div>
        )}
      </AnimatePresence>
      {status !== 'loading' && (
        <div className="mt-5 w-full max-w-md rounded-2xl border border-[#e8edf5] bg-white/90 p-4 text-left shadow-sm">
          <p className="text-sm font-semibold text-[#26344b]">Для жюри - демо-доступ</p>
          <p className="mt-1 text-xs leading-5 text-[#66758c]">
            Данные отделены от основной организации. Изменения сохраняются в демо-профиле.
          </p>
          <div className="mt-3 grid grid-cols-2 gap-2">
            <Button
              type="button"
              className="h-auto min-h-11 whitespace-normal px-2 py-2 text-sm"
              disabled={isDemoLoading}
              onClick={() => void enterDemo('teacher')}
            >
              Демо-преподаватель
            </Button>
            <Button
              type="button"
              className="h-auto min-h-11 whitespace-normal px-2 py-2 text-sm"
              disabled={isDemoLoading}
              onClick={() => void enterDemo('student')}
            >
              Демо-студент
            </Button>
          </div>
          {isDemoLoading && (
            <p role="status" className="mt-2 text-xs text-[#66758c]">
              Входим в демо-профиль…
            </p>
          )}
          {demoError && (
            <p role="alert" className="mt-2 text-sm text-red-600">
              Не удалось войти в демо-профиль. Попробуйте ещё раз позже.
            </p>
          )}
        </div>
      )}
    </motion.section>
  );
};
