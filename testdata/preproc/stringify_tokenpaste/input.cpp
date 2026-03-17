#define MOD_PREFIX Warpbox
#define MOD_ASSET_DIR assets/data
#define CFG_STR_INNER(X) #X
#define CFG_STR(X) CFG_STR_INNER(X)
#define ITEM_LOC_NAME_INNER(PREFIX, KEY) CFG_STR($##PREFIX##_##KEY##_name)
#define ITEM_LOC_NAME(PREFIX, KEY) ITEM_LOC_NAME_INNER(PREFIX, KEY)
#define MAT_DMG(BASE) CFG_STR(MOD_PREFIX/MOD_ASSET_DIR/BASE##_damage.rvmat)

class Demo
{
	displayName = ITEM_LOC_NAME(MOD_PREFIX, CardboardBox);
	damageMat = MAT_DMG(testbox);
};
