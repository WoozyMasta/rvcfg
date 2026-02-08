class CfgPatches
{
	class rvcfg_test_pants
	{
		requiredAddons[] = {};
	};
};

class cfgVehicles
{
	class Clothing;
	class TestPants_ColorBase: Clothing
	{
		scope = 2;
		displayName = "Test Pants";
		inventorySlot[] += {"Legs"};
	};

	class TestPants_Black: TestPants_ColorBase
	{
		hiddenSelectionsTextures[] = {"dz\\characters\\pants\\data\\testpants_black_co.paa"};
	};
};
