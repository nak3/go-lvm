#include <libdevmapper.h>
#include <lvm2app.h>
#include <stdio.h>

typedef char *names[];

struct result {
	char **namelist;
};

int wrapper_dm_list_iterate_items(struct dm_list *vgnames, char **r);
void wrapper_set_pv_prop(pv_create_params_t params, char *name, long long value);

char**makeCharArray(int size);

int is_valid(void* p);
int is_integer(void* p);
int is_signed(void* p);


