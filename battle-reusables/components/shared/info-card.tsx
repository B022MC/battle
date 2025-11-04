import React, { ReactNode } from 'react';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Text } from '@/components/ui/text';
import { View } from 'react-native';

type InfoCardProps = { children: ReactNode; className?: string };

export const InfoCard = ({ children, className = 'py-4 gap-4' }: InfoCardProps) => {
  return <Card className={className}>{children}</Card>;
};

type InfoCardHeaderProps = { children: ReactNode; className?: string };

export const InfoCardHeader = ({ children, className = 'px-4' }: InfoCardHeaderProps) => {
  return <CardHeader className={className}>{children}</CardHeader>;
};

type InfoCardTitleProps = { children: ReactNode; className?: string };

export const InfoCardTitle = ({ children, className }: InfoCardTitleProps) => {
  return <CardTitle className={className}>{children}</CardTitle>;
};

type InfoCardContentProps = { children: ReactNode; className?: string };

export const InfoCardContent = ({ children, className = 'px-4' }: InfoCardContentProps) => {
  return <CardContent className={className}>{children}</CardContent>;
};

type InfoCardRowProps = {
  label: string;
  value?: string | number;
  className?: string;
  classNames?: { label?: string; value?: string };
};

export const InfoCardRow = ({
  label,
  value,
  className = 'flex-row items-center justify-between',
  classNames,
}: InfoCardRowProps) => {
  return (
    <View className={className}>
      <Text variant="muted" className={classNames?.label}>
        {label}
      </Text>
      <Text variant="large" className={classNames?.value}>
        {value}
      </Text>
    </View>
  );
};

type InfoCardFooterProps = { children: ReactNode; className?: string };

export const InfoCardFooter = ({ children, className }: InfoCardFooterProps) => {
  return <CardFooter className={className}>{children}</CardFooter>;
};
