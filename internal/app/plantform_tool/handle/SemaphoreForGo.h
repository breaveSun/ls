#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <Windows.h>
//#include <iostream>

char addr[] = "up3d_new_task_notify";
HANDLE pSemaphore = NULL;

bool release()
{
	bool ret = ReleaseSemaphore(pSemaphore, 1, NULL);
	//std::cout << __FUNCTION__ << ret << std::endl;
	printf("%s %d \n", __FUNCTION__, ret);
	return ret;
};

bool notify()
{
	if (NULL == pSemaphore)
		pSemaphore = CreateSemaphoreA(NULL, 0, 1, addr);
	
	if (NULL == pSemaphore)
		return false;

	return release();
}

bool closeSemaphore()
{
	bool ret = CloseHandle(pSemaphore);
	//std::cout << __FUNCTION__ << ret << std::endl;
	printf("%s %d \n", __FUNCTION__, ret);
	pSemaphore = NULL;
	return ret;
};



	
