type UserData = {
  id: string;
  maxId: string;
};

export const userDataService = {
  setUserDataInLocalStorage(data: UserData) {
    let hasData = localStorage.getItem('_u-data');
    if (!hasData) localStorage.setItem('_u-data', JSON.stringify(data));
  },
  getUserDataFromLocalStorage(): UserData | null {
    let hasData = localStorage.getItem('_u-data');
    if (hasData) return JSON.parse(hasData) as UserData;
    return null;
  },
};
