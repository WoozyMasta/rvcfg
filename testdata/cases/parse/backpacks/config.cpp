class CfgPatches
{
  class MyMod_Characters_Backpacks
  {
    requiredAddons[] = {"DZ_Characters", "DZ_Characters_Backpacks"};
  };
};

class cfgVehicles
{
  class Clothing;

  class TaloonBag_ColorBase: Clothing
  {
    itemSize[] = {4, 5};
    itemsCargoSize[] = {5, 6};
    prohibitedInVehicles = 1;
    inventorySlot[] += {"Bag_1", "Bag_2", "Bag_3"};
  };

  class MountainBag_ColorBase: Clothing
  {
    itemSize[] = {5, 9};
    itemsCargoSize[] = {6, 10};
    attachments[] += {"CookingTripod"};
    prohibitedInVehicles = 1;
    inventorySlot[] += {"Bag_1", "Bag_2", "Bag_3"};
  };

  class ArmyPouch_ColorBase: Clothing
  {
    itemSize[] = {4, 4};
    itemsCargoSize[] = {4, 5};
    preventStackPouch = 1;
    inventorySlot[] += {"Hips", "Bag_1", "Bag_2"};
    itemInfo[] += {"Hips"};
  };

  class LargeTentBackPack: Clothing
  {
    scope = 2;
    itemSize[] = {4, 10};
    itemsCargoSize[] = {0, 0};
    inventorySlot[] += {"Bag_1", "Bag_2", "Bag_3"};

    class DamageSystem
    {
      class GlobalHealth
      {
        class Health
        {
          hitpoints = 2000;
          healthLevels[] =
          {
            {1.0, {"DZ\\gear\\camping\\data\\bagpack.rvmat"}},
            {0.5, {"DZ\\gear\\camping\\data\\bagpack_damage.rvmat"}},
            {0.0, {"DZ\\gear\\camping\\data\\bagpack_destruct.rvmat"}}
          };
        };
      };
    };
  };
};
