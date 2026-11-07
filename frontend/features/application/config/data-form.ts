import * as yup from 'yup';

export const applicationCreateRequestSchema = yup.object({
  name: yup
    .string()
    .min(10, 'Укажите полное ФИО')
    .max(50, 'Имя может содержать максимум 50 символов')
    .required('Введите имя!'),
  organizationId: yup.string().required('Выберите свою организацию!'),
  groupId: yup
    .string()
    .default('')
    .when('type', {
      is: 'student',
      then: (schema) => schema.required('Выберите свою группу!'),
      otherwise: (schema) => schema.defined(),
    }),
  userId: yup.string().required('Идентификатор пользователя не указан'),
  type: yup.string().required('Выберите тип заявки!'),
});
